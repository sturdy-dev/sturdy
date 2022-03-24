import { MutagenExecutable } from './Executable'
import { mkdir, readdir, stat } from 'fs/promises'
import path from 'path'
import { homedir, hostname } from 'os'
import { BrowserWindow, dialog } from 'electron'
import { createClient, gql } from '@urql/core'
import { SSHKeys } from '../SSHKeys'
import { MutagenSessionConfigurator } from './SessionConfigurator'
import { MutagenSession } from './Session'
import { MutagenIPC, sharedMutagenIpc } from '../ipc'
import { MutagenDaemon } from './Daemon'
import { dataPath } from '../resources'
import { File } from '../config'
import { Status, Auth } from '../application'
import { PostHogTracker } from '../PostHogTracker'
import { Logger } from '../Logger'

export class MutagenManager {
  readonly #logger: Logger
  readonly #mutagen: MutagenDaemon
  readonly #mutagenExecutable: MutagenExecutable
  readonly #status: Status
  readonly #configFile: File
  readonly #apiURL: URL
  readonly #graphqlURL: URL
  readonly #syncHostURL: URL
  readonly #postHog: PostHogTracker
  readonly #auth: Auth

  #sessions: MutagenSession[] = []
  #mainWindow: BrowserWindow | undefined

  constructor(
    logger: Logger,
    mutagen: MutagenDaemon,
    mutagenExecutable: MutagenExecutable,
    status: Status,
    configFile: File,
    apiURL: URL,
    graphqlURL: URL,
    syncHostURL: URL,
    postHog: PostHogTracker,
    auth: Auth
  ) {
    this.#logger = logger.withPrefix('mutagen-manager')
    this.#mutagen = mutagen
    this.#mutagenExecutable = mutagenExecutable
    this.#status = status
    this.#configFile = configFile
    this.#apiURL = apiURL
    this.#graphqlURL = graphqlURL
    this.#syncHostURL = syncHostURL
    this.#postHog = postHog
    this.#auth = auth

    mutagen.on('session-state-changed', (name, fromState, toState) => {
      // Wait for sessions to update their states
      setImmediate(() => {
        status.sessionsChanged(this.#sessions)
      })
    })

    auth.on('logged-in', () => {
      this.start()
    })

    auth.on('logged-out', () => {
      this.cleanup()
    })
  }

  setMainWindow(mainWindow: BrowserWindow) {
    this.#mainWindow = mainWindow
  }

  async preStartMutagen() {
    try {
      this.#logger.log('checking for cmd sturdy-client...')
      await stat(path.join(homedir(), '.sturdy-sync', 'daemon', 'daemon.sock'))
      dialog.showErrorBox(
        'Conflicting Sturdy Installations',
        "Looks like you're running the CLI version of Sturdy. The app and the CLI isn't supposed to work in tandem. Please run 'sturdy stop' and try to start the app again."
      )
      process.exit(1)
    } catch {}
  }

  start() {
    this.#start().catch((e: Error) => {
      this.#logger.error('failed to start mutagen, retrying once', e)

      // Try to force restart once, otherwise give up and let user decide
      // whether to force restart again.
      this.forceRestart().catch((e) => {
        this.#logger.error('failed to force restart mutagen', e)
      })
    })
  }

  async #start() {
    if (this.#mutagen.isRunning) {
      await this.cleanup()
    }

    const jwt = this.#auth.jwt

    if (!jwt) {
      return
    }

    const client = createClient({
      url: this.#graphqlURL.href,
      fetch: (await import('node-fetch')).default as any,
      fetchOptions: {
        credentials: 'include',
        headers: {
          Authorization: `bearer ${jwt}`,
        },
      },
    })

    const { data, error } = await client
      .query<{ user: { views: { id: string; mountPath: string }[] } }>(
        gql`
          {
            user {
              views {
                id
                mountPath
              }
            }
          }
        `
      )
      .toPromise()

    if (error != null || data?.user == null) {
      return
    }

    await this.#mutagen.start()

    const agentDir = dataPath(process.env.AGENT_DIR_NAME || 'sturdy-agent')
    await mkdir(agentDir, { recursive: true })

    const sshKeys = new SSHKeys(this.#logger, client, this.#status, this.#syncHostURL, agentDir)

    await sshKeys.ensure()

    const sessionConfigurator = new MutagenSessionConfigurator(
      this.#logger,
      agentDir,
      this.#mutagenExecutable,
      this.#mutagen,
      sshKeys,
      this.#apiURL,
      this.#syncHostURL,
      client
    )

    // Only the views that configured and exist on the server for this user are expected
    const expectedViews = this.#configFile.data.views.filter((configView) =>
      data.user.views.some(
        (apiView) => apiView.id === configView.id && apiView.mountPath === configView.path
      )
    )

    this.#sessions = await MutagenSession.reconcile(
      this.#logger,
      sessionConfigurator,
      expectedViews
    )

    this.#logger.log(
      'reconciled sessions:',
      this.#sessions.map((s) => s.name)
    )
    this.#status.reconciledSessions()

    const manager = this

    const logger = this.#logger

    const ipcImplementation: MutagenIPC = {
      async createView(workspaceID, mountPath) {
        logger.log(`createView ${workspaceID} at ${mountPath}`)

        const dirname = path.dirname(mountPath)
        const baseName = path.basename(mountPath)
        const content = await readdir(dirname)
        if (content.find((f) => f === baseName)) {
          throw new Error(`${baseName} already exists`)
        }

        const { error, data } = await client
          .mutation(
            gql`
              mutation ($input: CreateViewInput!) {
                createView(input: $input) {
                  id
                }
              }
            `,
            { input: { workspaceID, mountPath, mountHostname: hostname() } }
          )
          .toPromise()

        if (error != null) {
          throw error
        }

        const viewID = data?.createView?.id

        if (viewID == null) {
          throw new Error('Failed to create view')
        }

        await manager.#configFile.update((data) => {
          data.views.push({
            id: viewID,
            path: mountPath,
          })
        })

        manager.#postHog.trackCreateNewView(viewID)

        manager.#sessions ??= []
        manager.#sessions.push(await sessionConfigurator.configureAndStart(viewID, mountPath))

        return viewID
      },
      async createNewViewWithDialog(workspaceID: string, codebaseSlug: string) {
        const { canceled, filePath } = await dialog.showSaveDialog(manager.#mainWindow!, {
          title: 'Select location',
          defaultPath: path.join(homedir(), codebaseSlug),
          buttonLabel: 'Select',
          nameFieldLabel: 'Open As:',
          showsTagField: false,
          properties: ['createDirectory', 'showOverwriteConfirmation'],
        })

        if (canceled || !filePath) {
          throw new Error('Cancelled by user')
        }

        try {
          return await this.createView(workspaceID, filePath)
        } catch (e) {
          logger.error(`failed to create view for workspace ${workspaceID} at ${filePath}`, e)
          throw e
        }
      },
    }

    Object.values(sharedMutagenIpc).forEach((method) => method.clean())
    Object.entries(ipcImplementation).forEach(([channel, implementation]) => {
      sharedMutagenIpc[channel as keyof MutagenIPC].implement(
        implementation.bind(ipcImplementation) as any
      )
    })
  }

  async cleanup() {
    try {
      await Promise.all(
        this.#sessions?.map((session) => session.pause().catch((e) => this.#logger.error(e))) ?? []
      )
    } catch (e) {
      this.#logger.error(e)
    }

    try {
      await this.#mutagen.stop()
    } catch (e) {
      this.#logger.error(e)
    }
  }

  async forceRestart() {
    await this.#mutagen.kill()
    this.#sessions = []
    await this.#mutagen.deleteDir()
    await this.#start()
  }
}

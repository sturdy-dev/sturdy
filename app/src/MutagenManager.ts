import { MutagenExecutable } from './MutagenExecutable'
import { mkdir, readdir, stat } from 'fs/promises'
import path from 'path'
import { homedir, hostname } from 'os'
import { BrowserWindow, dialog } from 'electron'
import { createClient, gql } from '@urql/core'
import { SSHKeys } from './SSHKeys'
import { MutagenSessionConfigurator } from './MutagenSessionConfigurator'
import { MutagenSession } from './MutagenSession'
import { MutagenIPC, sharedMutagenIpc } from './ipc'
import { MutagenDaemon } from './MutagenDaemon'
import { dataPath } from './resources'
import { ConfigFile } from './ConfigFile'
import { AppStatus } from './AppStatus'
import { PostHogTracker } from './PostHogTracker'
import { Auth } from './Auth'

export class MutagenManager {
  readonly #mutagen: MutagenDaemon
  readonly #mutagenExecutable: MutagenExecutable
  readonly #status: AppStatus
  #sessions: MutagenSession[] = []
  readonly #configFile: ConfigFile
  readonly #apiURL: URL
  readonly #syncHostURL: URL
  readonly #postHog: PostHogTracker
  readonly #auth: Auth

  #mainWindow: BrowserWindow | undefined

  constructor(
    mutagen: MutagenDaemon,
    mutagenExecutable: MutagenExecutable,
    status: AppStatus,
    configFile: ConfigFile,
    apiURL: URL,
    syncHostURL: URL,
    postHog: PostHogTracker,
    auth: Auth
  ) {
    this.#mutagen = mutagen
    this.#mutagenExecutable = mutagenExecutable
    this.#status = status
    this.#configFile = configFile
    this.#apiURL = apiURL
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
      await stat(path.join(homedir(), '.sturdy-sync', 'daemon', 'daemon.sock'))
      await dialog.showErrorBox(
        'Conflicting Sturdy Installations',
        "Looks like you're running the CLI version of Sturdy. The app and the CLI isn't supposed to work in tandem. Please run 'sturdy stop' and try to start the app again."
      )
      process.exit(1)
      return
    } catch {}
  }

  start() {
    this.#start().catch((e: Error) => {
      console.log('failed to start mutagen, retrying once', e)

      // Try to force restart once, otherwise give up and let user decide
      // whether to force restart again.
      this.forceRestart().catch((e) => {
        console.log('failed to force restart mutagen', e)
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
      url: this.#apiURL.href,
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

    const sshKeys = new SSHKeys(client, this.#status, this.#syncHostURL, agentDir)

    await sshKeys.ensure()

    const sessionConfigurator = new MutagenSessionConfigurator(
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

    this.#sessions = await MutagenSession.reconcile(sessionConfigurator, expectedViews)

    console.log('Reconciled sessions:', this.#sessions)
    this.#status.reconciledSessions()

    let manager = this

    const ipcImplementation: MutagenIPC = {
      async createView(workspaceID, mountPath) {
        if ((await readdir(mountPath)).length > 0) {
          throw new Error('Cannot create view in non-empty directory')
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
      async createNewViewWithDialog(workspaceID: string) {
        const { canceled, filePaths } = await dialog.showOpenDialog(manager.#mainWindow!, {
          properties: ['openDirectory', 'createDirectory'],
        })
        if (canceled || filePaths.length === 0) {
          throw new Error('Cancelled by user')
        }

        return this.createView(workspaceID, filePaths[0])
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
        this.#sessions?.map((session) => session.pause().catch((e) => console.error(e))) ?? []
      )
    } catch (e) {
      console.error(e)
    }

    try {
      await this.#mutagen.stop()
    } catch (e) {
      console.error(e)
    }
  }

  async forceRestart() {
    await this.#mutagen.kill()
    this.#sessions = []
    await this.#mutagen.deleteDir()
    await this.#start()
  }
}

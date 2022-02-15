import { MutagenSessionConfigurator } from './SessionConfigurator'
import { ViewConfig } from '../Config'
import { MutagenExecutable } from './Executable'
import { MutagenDaemon } from './Daemon'
import { Logger } from '../Logger'

// SESSION_VERSION_NUMBER tracks versions of the sturdy configuration.
//
// Sessions are labeled with this number
// If the existing label is different from this number, the session will be terminated and re-created.
// 5 - Sturdy < 0.2.5
// [6-19] (inclusive) - Are reserved for the CLI version of Sturdy
// 20 - Sturdy > 0.3.0
// 21 - Strudy >= 0.5.0
const SESSION_VERSION_NUMBER = 21

interface ListOutputSession {
  identifier: string
  name: string
  alpha: {
    path: string
  }
  labels: Record<string, string>
  // Add more here as needed
}

export interface SessionSharedConfig {
  logger: Logger
  configPath: string
  userID: string
  apiURL: URL
  syncHostURL: URL
  executable: MutagenExecutable
  daemon: MutagenDaemon
}

export interface SessionConfig extends SessionSharedConfig {
  viewID: string
  codebaseID: string
  mountPath: string
}

export enum MutagenSessionState {
  Disconnected = 'Disconnected',
  HaltedOnRootEmptied = 'HaltedOnRootEmptied',
  HaltedOnRootDeletion = 'HaltedOnRootDeletion',
  HaltedOnRootTypeChange = 'HaltedOnRootTypeChange',
  ConnectingAlpha = 'ConnectingAlpha',
  ConnectingBeta = 'ConnectingBeta',
  Watching = 'Watching',
  Scanning = 'Scanning',
  WaitingForRescan = 'WaitingForRescan',
  Reconciling = 'Reconciling',
  StagingAlpha = 'StagingAlpha',
  StagingBeta = 'StagingBeta',
  Transitioning = 'Transitioning',
  Saving = 'Saving',
  Unknown = 'Unknown',
}

export class MutagenSession {
  readonly #logger: Logger
  readonly #executable: MutagenExecutable
  readonly #daemon: MutagenDaemon
  readonly #sessionVersion: string
  readonly name: string
  readonly viewID: string
  readonly path: string
  #state = MutagenSessionState.Unknown

  #unsubscribeToStateUpdates!: () => void
  #isSubscribedToStateUpdates: boolean = false

  constructor(
    logger: Logger,
    executable: MutagenExecutable,
    daemon: MutagenDaemon,
    sessionVersion: string,
    name: string,
    viewID: string,
    path: string
  ) {
    this.#logger = logger.withPrefix('session', name)
    this.#executable = executable
    this.#daemon = daemon
    this.#sessionVersion = sessionVersion
    this.name = name
    this.viewID = viewID
    this.path = path

    this.#subscribeToStateUpdates()
  }

  #subscribeToStateUpdates() {
    if (this.#isSubscribedToStateUpdates) {
      return
    }

    this.#logger.log('subscribe to status updates')

    let listener = (name: string, oldState: MutagenSessionState, newState: MutagenSessionState) => {
      if (name === this.name) {
        this.#state = newState
      }
    }

    this.#daemon.on('session-state-changed', listener)
    this.#isSubscribedToStateUpdates = true
    this.#unsubscribeToStateUpdates = () => {
      this.#daemon.off('session-state-changed', listener)
      this.#isSubscribedToStateUpdates = false
    }
  }

  get getState(): string {
    return this.#state
  }

  get isRunning(): boolean {
    switch (this.#state) {
      case MutagenSessionState.Disconnected:
      case MutagenSessionState.ConnectingAlpha:
      case MutagenSessionState.ConnectingBeta:
      case MutagenSessionState.HaltedOnRootEmptied:
      case MutagenSessionState.HaltedOnRootDeletion:
      case MutagenSessionState.HaltedOnRootTypeChange:
      case MutagenSessionState.Unknown:
        return false
      case MutagenSessionState.Watching:
      case MutagenSessionState.Scanning:
      case MutagenSessionState.Reconciling:
      case MutagenSessionState.Transitioning:
      case MutagenSessionState.Saving:
      case MutagenSessionState.WaitingForRescan:
      case MutagenSessionState.StagingAlpha:
      case MutagenSessionState.StagingBeta:
        return true
      default:
        this.#logger.error('unexpected mutagen state', this.#state)
        return true
    }
  }

  static async reconcile(
    logger: Logger,
    configurator: MutagenSessionConfigurator,
    expectedViews: ViewConfig[]
  ): Promise<MutagenSession[]> {
    const existingSessions = await MutagenSession.list(
      logger,
      configurator.executable,
      configurator.daemon,
      configurator.apiURL
    )

    // Views that are configured but doesn't have a session
    const missingViews = expectedViews.filter(
      (view) => !existingSessions.some((session) => session.viewID === view.id)
    )

    // Sessions that aren't configured but exists
    const unexpectedSessions = existingSessions.filter(
      (session) => !expectedViews.some((view) => session.viewID === view.id)
    )

    // Sessions that are configured and exists already
    const expectedSessions = existingSessions.filter((session) =>
      expectedViews.some((view) => session.viewID === view.id)
    )

    const [sessions] = await Promise.all([
      // Resuming/creating sessions
      Promise.all([
        ...expectedSessions.map(async (s) => {
          if (s.isStale) {
            await s.terminate()
            return configurator.configureAndStart(s.viewID, s.path)
          }
          await s.resume()
          return [s]
        }),
        ...missingViews.map(async (v) => {
          try {
            const session = await configurator.configureAndStart(v.id, v.path)
            return [session]
          } catch (e) {
            logger.error(e)
            return []
          }
        }),
      ]),

      // Side effects
      Promise.all([...unexpectedSessions.map((s) => s.terminate())]),
    ])

    return sessions.flat()
  }

  static async list(
    logger: Logger,
    executable: MutagenExecutable,
    daemon: MutagenDaemon,
    apiURL: URL
  ): Promise<MutagenSession[]> {
    const [list, onExit] = executable.execute(['sync', 'list', '--json'], {
      stdio: ['ignore', 'pipe', 'ignore'],
    })

    const chunks: Buffer[] = []
    list.stdout.on('data', (chunk) => chunks.push(chunk))

    await onExit

    const json = Buffer.concat(chunks).toString('utf-8')

    const sessions: null | { session: ListOutputSession }[] = JSON.parse(json)

    return (
      sessions
        ?.map((s) => s.session)
        .filter((s) => s.labels?.sturdy === 'true')
        .filter((s) => s.labels?.sturdyApiHost === apiURL.hostname)
        .map((s) => {
          return new MutagenSession(
            logger,
            executable,
            daemon,
            s.labels.sessionVersion,
            s.name,
            s.name.slice('view-'.length),
            s.alpha.path
          )
        }) ?? []
    )
  }

  static async create({
    logger,
    configPath,
    mountPath,
    userID,
    viewID,
    codebaseID,
    apiURL,
    syncHostURL,
    executable,
    daemon,
  }: SessionConfig): Promise<MutagenSession> {
    const name = `view-${viewID}`
    logger.log(`new session ${name}`)
    // prettier-ignore
    const args = [
            'sync', 'create',
            '--no-global-configuration',
            '-c', configPath,
            '--name', name,
            '--label', 'sturdy=true',
            '--label', `sessionVersion=${SESSION_VERSION_NUMBER}`,
            '--label', `sturdyApiProto=${apiURL.protocol.slice(0, -1)}`,
            '--label', `sturdyApiHost=${apiURL.hostname}`,
            '--label', `sturdyApiHostPort=${apiURL.port}`,
            '--label', `sturdyApiPrefix=${apiURL.pathname === '' ? '' : apiURL.pathname.slice(1)}`,
            '--label', `sturdyViewId=${viewID}`,
            '--stage-mode-beta=neighboring',

            // Alpha
            mountPath,

            // Beta
            // Note: /repos here is a convention
            `${userID}@${syncHostURL.host}:/repos/${codebaseID}/${viewID}/`,
        ]

    // Create session object before calling the daemon to create it
    // This allows us to catch events from during the creation process _inside_ the MutagenSession
    let newSession = new MutagenSession(
      logger,
      executable,
      daemon,
      SESSION_VERSION_NUMBER.toString(),
      name,
      viewID,
      mountPath
    )

    const [, onExit] = executable.execute(args, {
      stdio: ['ignore', 'ignore', 'ignore'],
    })
    await onExit

    return newSession
  }

  async pause() {
    const [, onExit] = this.#executable.execute(['sync', 'pause', this.name], {
      stdio: ['ignore', 'ignore', 'ignore'],
      timeout: 10000,
    })
    await onExit
    this.#unsubscribeToStateUpdates()
  }

  async resume() {
    this.#subscribeToStateUpdates()
    const [, onExit] = this.#executable.execute(['sync', 'resume', this.name], {
      stdio: ['ignore', 'ignore', 'ignore'],
    })
    await onExit
  }

  get isStale() {
    return this.#sessionVersion !== SESSION_VERSION_NUMBER.toString()
  }

  async terminate() {
    const [, onExit] = this.#executable.execute(['sync', 'terminate', this.name], {
      stdio: ['ignore', 'ignore', 'ignore'],
    })
    await onExit
    this.#unsubscribeToStateUpdates()
  }
}

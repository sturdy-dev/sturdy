import path from 'path'
import { Host } from './Host'
import { PostHogTracker } from '../PostHogTracker'
import { Auth } from './Auth'
import { MutagenManager, MutagenExecutable, MutagenDaemon } from '../mutagen'
import { dataPath, resourcePath } from '../resources'
import { Status } from './Status'
import { dialog, BrowserWindow, shell, app } from 'electron'
import { Logger } from '../Logger'
import { File } from '../config'
import { AppIPC, MutagenIPC, sharedAppIpc, sharedMutagenIpc } from '../ipc'
import { TypedEventEmitter } from '../TypedEventEmitter'

interface ApplicationEvents {
  openPreferences: []
}

export class Application extends TypedEventEmitter<ApplicationEvents> {
  readonly #host: Host
  readonly #auth: Auth
  readonly #isAppPackaged: boolean
  readonly #protocol: string
  readonly #mutagenManager: MutagenManager
  readonly #postHogTracker: PostHogTracker
  readonly #status: Status
  readonly #logger: Logger

  #window?: BrowserWindow
  #lastURL?: URL

  private constructor(
    host: Host,
    auth: Auth,
    isAppPackaged: boolean,
    protocol: string,
    mutagenManager: MutagenManager,
    postHogTracker: PostHogTracker,
    status: Status,
    logger: Logger
  ) {
    super()
    this.#host = host
    this.#auth = auth
    this.#isAppPackaged = isAppPackaged
    this.#protocol = protocol
    this.#postHogTracker = postHogTracker
    this.#mutagenManager = mutagenManager
    this.#status = status
    this.#logger = logger
  }

  static async start({
    host,
    mutagenExecutable,
    postHogToken,
    isAppPackaged,
    protocol,
    logger,
    daemon,
    status,
  }: {
    host: Host
    mutagenExecutable: MutagenExecutable
    postHogToken: string
    isAppPackaged: boolean
    protocol: string
    logger: Logger
    daemon: MutagenDaemon
    status: Status
  }) {
    logger = logger.withPrefix(host.id)
    logger.log('starting')

    const auth = await Auth.start(host.graphqlURL)

    const configFile =
      host.id === 'cloud'
        ? await File.open(dataPath('config.json')) // fallback to support apps before the multi-backend update
        : await File.open(dataPath(host.id, 'config.json'))
    logger.log('config file', configFile.path)

    const postHogTracker = new PostHogTracker(host.graphqlURL, postHogToken)
    const mutagenManager = new MutagenManager(
      logger,
      daemon,
      mutagenExecutable,
      status,
      configFile,
      host.apiURL,
      host.graphqlURL,
      host.syncURL,
      postHogTracker,
      auth
    )

    await mutagenManager.preStartMutagen()

    const migrated = await configFile.migrate()
    if (migrated) {
      await dialog.showMessageBox({
        title: 'Migration Complete',
        message:
          "Thanks for being a Sturdy user! We've migrated the configuration from your existing Sturdy installation to the native app. Any questions? Reach out to support@getsturdy.com!",
      })
    }

    auth.on('logged-in', async () => {
      logger.log('logged-in')

      if (auth.jwt) {
        await postHogTracker
          .updateUser(auth.jwt)
          .catch((e) => logger.error('failed to setup posthog tracker', { e }))
        postHogTracker.trackStartedApp()
      }
    })

    auth.on('logged-out', async () => {
      logger.log('logged-out')

      postHogTracker.unsetUser()
    })

    if (auth.jwt != null) {
      logger.log('starting with existing jwt')

      await postHogTracker
        .updateUser(auth.jwt)
        .catch((e) => logger.error('failed to setup posthog tracker', { e }))
      postHogTracker.trackStartedApp()

      mutagenManager.start()
    }

    return new Application(
      host,
      auth,
      isAppPackaged,
      protocol,
      mutagenManager,
      postHogTracker,
      status,
      logger
    )
  }

  async close() {
    if (this.#window) {
      this.#logger.log('closing window')
      this.#window.close()
    }
  }

  async open(startURL?: URL) {
    if (!startURL && this.#lastURL) {
      startURL = this.#lastURL
    }

    // Re-use window if exists
    if (this.#window != null) {
      this.#logger.log('re-using existing window')
      if (this.#window.isMinimized()) {
        this.#window.restore()
      }
      this.#window.show()
      this.#window.focus()
      if (startURL != null) {
        this.#logger.log('opening url in existing window', startURL.href)
        await loadURLWithoutThrowingOnRedirects(this.#window, this.#logger, startURL.href)
      }
      return this.#window
    }

    this.#logger.log('creating new window')
    app.dock?.show()

    this.#window = new BrowserWindow({
      height: 1200,
      width: 1800,
      minWidth: 680,
      minHeight: 400,
      webPreferences: {
        spellcheck: true,
        preload: path.join(__dirname, 'preload.js'),
        devTools: !this.#isAppPackaged,
      },
      titleBarStyle: 'hidden',
      titleBarOverlay: {
        color: '#F9FAFB',
        symbolColor: '#1F2937',
      },
      trafficLightPosition: { x: 16, y: 16 },
    })

    // Create base IPC
    this.#addFallbackMutagenIpc()
    this.#addNonMutagenIpc()

    this.#window.removeMenu()

    const url = new URL(startURL ?? this.#host.webURL.href)
    if (startURL == null) {
      if (this.#auth.jwt == null) {
        url.pathname = '/login'
      } else {
        url.pathname = '/codebases'

        for (const arg of process.argv) {
          if (arg.startsWith(this.#protocol + '://')) {
            try {
              url.pathname = new URL(arg).pathname
            } catch (e) {
              this.#logger.error(e)
            }
          }
        }
      }
    }

    this.#window.once('closed', () => {
      this.#window = undefined
    })

    this.#window.webContents.on('did-navigate-in-page', (_, url) => {
      const targetURL = new URL(url)
      if (targetURL.toString() === this.#lastURL?.toString()) return
      this.#logger.log('navigated to', targetURL.toString())
      this.#lastURL = new URL(targetURL)
    })

    this.#window.webContents.on('before-input-event', (event, input) => {
      const cmdOrCtrl = input.meta || input.control
      if (cmdOrCtrl && input.key.toLowerCase() === ',') {
        event.preventDefault()
        this.emit('openPreferences')
      }
    })

    this.#window.webContents.on('will-navigate', (event, url) => {
      if (!url.startsWith(this.#host.webURL.href)) {
        shell.openExternal(url)
        event.preventDefault()
      }
    })

    // open target="_blank" links eternally
    this.#window.webContents.setWindowOpenHandler(({ url }) => {
      shell.openExternal(url)
      return { action: 'deny' }
    })

    // If the page fails to load, display app-fail.html
    this.#window.webContents.on(
      'did-fail-load',
      (
        event: Event,
        errorCode: number,
        errorDescription: string,
        validatedURL: string,
        isMainFrame: boolean,
        frameProcessId: number,
        frameRoutingId: number
      ) => {
        this.#logger.error('did-fail-load', { errorCode, errorDescription })
        this.#window?.loadFile(resourcePath('app-fail.html'), {
          query: {
            goto: this.#host.webURL.href,
          },
        })
      }
    )

    this.#mutagenManager.setMainWindow(this.#window)

    try {
      this.#logger.log('loading url', url.href)
      await loadURLWithoutThrowingOnRedirects(this.#window, this.#logger, url.href)
    } catch (e) {
      this.#logger.error('failed to loadURL', e)
    }
  }

  async isUp() {
    return await this.#host.isUp()
  }

  get status() {
    return this.#status
  }

  async cleanup() {
    this.#logger.log('cleaning up...')
    this.#postHogTracker.flush()
    await this.#mutagenManager?.cleanup()
  }

  get host() {
    return this.#host
  }

  async forceRestart() {
    try {
      await this.#mutagenManager.forceRestart()
    } catch (e) {
      this.#logger.error('failed to restart mutagen', e)
    }
  }

  #addFallbackMutagenIpc() {
    const ipcImplementation: MutagenIPC = {
      async createView(workspaceID, mountPath) {
        throw new Error('mutagen is not available')
      },
      async createNewViewWithDialog(workspaceID: string) {
        throw new Error('mutagen is not available')
      },
    }

    Object.values(sharedMutagenIpc).forEach((method) => method.clean())
    Object.entries(ipcImplementation).forEach(([channel, implementation]) => {
      sharedMutagenIpc[channel as keyof MutagenIPC].implement(
        implementation.bind(ipcImplementation) as any
      )
    })
  }

  #addNonMutagenIpc() {
    const auth = this.#auth
    const window = this.#window
    const logger = this.#logger
    const mutagenManager = this.#mutagenManager
    const status = this.#status

    const ipcImplementation: AppIPC = {
      isAuthenticated() {
        return auth.jwt !== null
      },
      goBack() {
        window?.webContents.goBack()
      },
      goForward() {
        window?.webContents.goForward()
      },
      canGoBack() {
        return window?.webContents.canGoBack() ?? false
      },
      canGoForward() {
        return window?.webContents.canGoForward() ?? false
      },
      state() {
        return status.state
      },
      async forceRestartMutagen() {
        try {
          await mutagenManager.forceRestart()
        } catch (e) {
          logger.error('failed to restart mutagen', e)
        }
      },
    }

    Object.entries(ipcImplementation).forEach(([channel, implementation]) => {
      sharedAppIpc[channel as keyof AppIPC].implement(implementation.bind(ipcImplementation) as any)
    })
  }
}

async function loadURLWithoutThrowingOnRedirects(
  window: BrowserWindow,
  logger: Logger,
  url: string
) {
  const newURL = new URL(url, window.webContents.getURL() || undefined)
  try {
    await window.loadURL(newURL.href)
  } catch (e) {
    if (typeof e === 'object' && e && 'code' in e && (e as any).code === 'ERR_ABORTED') {
      // This error is emitted if the browser redirects immediately after
      // loading the requested URL, which happens in the SPA for different
      // reasons. So we don't want this to become an actual error.

      // The exception is if the browser has actually remained on the
      // page that we navigated to, but still produced the error, because
      // then there was something that actually aborted the navigation.
      if (window.webContents.getURL() === newURL.href) {
        throw e
      }

      logger.log(
        'caught redirected loadURL from',
        newURL.href,
        'to',
        window.webContents.getURL(),
        e
      )
      return
    }
    throw e
  }
}

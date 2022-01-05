import path from 'path'
import {
  app,
  BrowserWindow,
  crashReporter,
  dialog,
  Menu,
  MenuItem,
  nativeImage,
  shell,
  Tray,
} from 'electron'
import { AppIPC, MutagenIPC, sharedAppIpc, sharedMutagenIpc } from './ipc'
import { Auth } from './Auth'
import { ConfigFile } from './ConfigFile'
import { MutagenExecutable } from './MutagenExecutable'
import { Updater } from './Updater'
import { MutagenManager } from './MutagenManager'
import { dataPath, resourceBinary, resourcePath } from './resources'
import { MutagenDaemon } from './MutagenDaemon'
import * as Sentry from '@sentry/electron'
import { AppStatus } from './AppStatus'
import log from 'electron-log'
import { createWriteStream, mkdirSync } from 'fs'
import { PostHogTracker } from './PostHogTracker'
import { CaptureConsole } from '@sentry/integrations'

// Start crash reporter before setting up logging
crashReporter.start({
  companyName: 'Sturdy Sweden AB',
  productName: 'Sturdy',
  ignoreSystemCrashHandler: true,
  submitURL:
    'https://o952367.ingest.sentry.io/api/6075838/minidump/?sentry_key=59a9e2de840941b58b49f82b0732e170',
})

const logsDir = dataPath('logs')
mkdirSync(logsDir, { recursive: true })

// Setup logging to file after crash reporter.
Object.assign(console, log.functions)
log.transports.file.resolvePath = () => path.join(logsDir, 'main.log')

let mainWindow: BrowserWindow | undefined

if (!app.requestSingleInstanceLock()) {
  app.quit()
} else {
  // Setup error logging
  // https://sentry.io/organizations/sturdy-xd/projects/sturdy-electron
  if (app.isPackaged) {
    Sentry.init({
      dsn: 'https://59a9e2de840941b58b49f82b0732e170@o952367.ingest.sentry.io/6075838',
      // release: "Sturdy@" + process.env.npm_package_version,
      // environment: "production",
      sampleRate: 1.0,
      integrations: [
        new CaptureConsole({
          levels: ['error'],
        }),
      ],
    })
  }

  const protocol = process.env.STURDY_PROTOCOL ?? 'sturdy'

  if (!app.isPackaged) {
    if (process.argv.length >= 2) {
      app.setAsDefaultProtocolClient(protocol, process.execPath, [path.resolve(process.argv[1])])
    }
  } else {
    app.setAsDefaultProtocolClient(protocol)
  }

  let urlProvidedOnStartup: string | undefined
  app.on('will-finish-launching', () =>
    app.on('open-url', (_event, url) => (urlProvidedOnStartup = url))
  )

  const iconSm = nativeImage.createFromPath(resourcePath('AppIconSm.png'))
  const iconSmDisconnected = nativeImage.createFromPath(resourcePath('AppIconSmDisconnected.png'))
  const iconSmTemplate = nativeImage.createFromPath(resourcePath('AppIconSmTemplate.png'))
  const iconSmDisconnectedTemplate = nativeImage.createFromPath(
    resourcePath('AppIconSmDisconnectedTemplate.png')
  )

  const iconTray = process.platform === 'darwin' ? iconSmTemplate : iconSm
  const iconTrayDisconnected =
    process.platform === 'darwin' ? iconSmDisconnectedTemplate : iconSmDisconnected

  const iconLg = nativeImage.createFromPath(resourcePath('AppIconLg.png'))
  const iconXl = nativeImage.createFromPath(resourcePath('AppIconXL.png'))

  const webURL = new URL(process.env.STURDY_WEB_URL ?? 'https://getsturdy.com')
  const apiURL = new URL(process.env.STURDY_API_URL ?? 'https://api.getsturdy.com/graphql')
  const syncHostURL = new URL(process.env.STURDY_SYNC_URL ?? 'ssh://sync.getsturdy.com')
  const postHogToken =
    process.env.STURDY_POSTHOG_API_KEY ?? 'ZuDRoGX9PgxGAZqY4RF9CCJJLpx14h3szUPzm7XBWSg'

  const runAutoUpdater = app.isPackaged && !process.env.STURDY_DISABLE_AUTO_UPDATER

  const longLivedInstances = new Set()

  // Save on global to avoid GC
  ;(global as any)[Symbol('LONG_LIVED_INSTANCES')] = longLivedInstances

  const status = new AppStatus()

  const mutagenLog = createWriteStream(path.join(logsDir, 'mutagen.log'), {
    flags: 'a',
    mode: 0o666,
  })

  const mutagenDataDirectoryPath = dataPath(process.env.MUTAGEN_DIR_NAME || 'mutagen')

  const mutagenExecutable = new MutagenExecutable({
    executablePath: resourceBinary('sturdy-sync'),
    dataDirectory: mutagenDataDirectoryPath,
    log: mutagenLog,
  })

  let mutagenManager: MutagenManager

  let postHog = new PostHogTracker(apiURL, postHogToken)

  async function main() {
    if (app.isPackaged) {
      await Updater.finalizePendingUpdate()
    }

    await app.whenReady()

    // Adds support for notifications on Windows
    if (process.platform === 'win32') {
      app.setAppUserModelId('com.getsturdy.sturdy')
    }

    const auth = await Auth.start(apiURL)

    const configFile = await ConfigFile.open(dataPath('config.json'))
    const daemon = new MutagenDaemon(mutagenExecutable, mutagenDataDirectoryPath, mutagenLog)
    mutagenManager = new MutagenManager(
      daemon,
      mutagenExecutable,
      status,
      configFile,
      apiURL,
      syncHostURL,
      postHog,
      auth
    )

    daemon.on('session-manager-initialized', status.mutagenStarted.bind(status))

    daemon.on('failed-to-start', status.mutagenFailedToStart.bind(status))
    daemon.on('is-running-changed', (isRunning) => {
      if (!isRunning) {
        status.mutagenStopped()
      }
    })

    daemon.on('connection-to-server-dropped', status.connectionDropped.bind(status))

    const openMenuItem = new MenuItem({
      label: 'Open Sturdy',
      click: () => triggerWindow(),
    })
    const quitMenuItem = new MenuItem({
      label: 'Quit Sturdy',
      click: quit,
      accelerator: 'CommandOrControl+Q',
    })

    const forceRestartMutagenMenuItem = new MenuItem({
      label: 'Force Restart Syncer',
      click: async () => {
        try {
          await mutagenManager.forceRestart()
        } catch (e) {
          console.error('failed to restart mutagen', e)
        }
      },
    })

    const debugMenuItem = new MenuItem({
      label: 'Debug',
      submenu: Menu.buildFromTemplate([forceRestartMutagenMenuItem]),
    })

    const tray = new Tray(iconTrayDisconnected)

    status.on('change', (state) => {
      if (state === 'online') {
        tray.setImage(iconTray)
      } else {
        tray.setImage(iconTrayDisconnected)
      }
    })

    const contextMenu = new Menu()

    status.appendMenuItem(contextMenu)

    contextMenu.append(new MenuItem({ type: 'separator' }))

    contextMenu.append(openMenuItem)
    contextMenu.append(debugMenuItem)

    contextMenu.append(new MenuItem({ type: 'separator' }))

    if (runAutoUpdater) {
      Updater.start(contextMenu).catch((e) => {
        console.error('Auto updater error!', e)
      })
    }

    contextMenu.append(quitMenuItem)

    tray.setContextMenu(contextMenu)

    longLivedInstances.add(tray)

    const migrated = await configFile.migrate()
    if (migrated) {
      await dialog.showMessageBox({
        title: 'Migration Complete',
        message:
          "Thanks for being a Sturdy user! We've migrated the configuration from your existing Sturdy installation to the native app. Any questions? Reach out to support@getsturdy.com!",
      })
    }

    // Before starting the main window, check if the user has a legacy (non-app) sturdy-sync mutagen session running
    console.log('Running preStartMutagen')
    await mutagenManager.preStartMutagen()
    console.log('Done preStartMutagen')

    function addFallbackMutagenIpc() {
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

    function addNonMutagenIpc() {
      const ipcImplementation: AppIPC = {
        isAuthenticated() {
          return auth.jwt != null
        },
        goBack() {
          mainWindow?.webContents.goBack()
        },
        goForward() {
          mainWindow?.webContents.goForward()
        },
        canGoBack() {
          return mainWindow?.webContents.canGoBack() ?? false
        },
        canGoForward() {
          return mainWindow?.webContents.canGoForward() ?? false
        },
        state() {
          return status.state
        },
        async forceRestartMutagen() {
          await mutagenManager.forceRestart()
        },
      }

      Object.entries(ipcImplementation).forEach(([channel, implementation]) => {
        sharedAppIpc[channel as keyof AppIPC].implement(
          implementation.bind(ipcImplementation) as any
        )
      })
    }

    // Create base IPC
    addFallbackMutagenIpc()
    addNonMutagenIpc()

    auth.on('logged-in', async () => {
      console.log('logged-in')

      if (auth.jwt) {
        await postHog
          .updateUser(auth.jwt)
          .catch((e) => console.error('failed to setup posthog tracker', { e }))
        postHog.trackStartedApp()
      }
    })

    auth.on('logged-out', async () => {
      console.log('logged-out')

      postHog.unsetUser()
    })

    if (auth.jwt != null) {
      console.log('starting with existing jwt')

      await postHog
        .updateUser(auth.jwt)
        .catch((e) => console.error('failed to setup posthog tracker', { e }))
      postHog.trackStartedApp()

      mutagenManager.start()
    }

    app.on('before-quit', (event) => {
      event.preventDefault()
      postHog.flush()
      quit()
    })

    app.on('window-all-closed', () => {
      // Don't do anything
      // Keep the app running in the tray
    })

    async function loadURLWithoutThrowingOnRedirects(window: BrowserWindow, url: string) {
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

          console.log(
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

    async function triggerWindow(withUrl?: URL) {
      console.log('triggerWindow', mainWindow != null, withUrl)

      // Re-use window if exists
      if (mainWindow != null) {
        if (mainWindow.isMinimized()) {
          mainWindow.restore()
        }
        mainWindow.show()
        mainWindow.focus()
        if (withUrl != null) {
          await loadURLWithoutThrowingOnRedirects(mainWindow, withUrl.href)
        }
        return
      }

      openMenuItem.label = 'Show Sturdy'
      app.dock?.show()

      mainWindow = new BrowserWindow({
        height: 1200,
        width: 1800,
        minWidth: 680,
        minHeight: 400,
        webPreferences: {
          preload: path.join(__dirname, 'preload.js'),
          devTools: !app.isPackaged,
        },
        titleBarStyle: 'hidden',
        titleBarOverlay: {
          color: '#F9FAFB',
          symbolColor: '#1F2937',
        },
        trafficLightPosition: { x: 16, y: 16 },
      })

      mainWindow.removeMenu()

      const url = new URL(withUrl ?? webURL.href)
      if (withUrl == null) {
        if (auth.jwt == null) {
          url.pathname = '/login'
        } else {
          url.pathname = '/codebases'

          for (const arg of process.argv) {
            if (arg.startsWith(protocol + '://')) {
              try {
                url.pathname = new URL(arg).pathname
              } catch (e) {
                console.error(e)
              }
            }
          }

          if (urlProvidedOnStartup != null) {
            try {
              url.pathname = new URL(urlProvidedOnStartup).pathname
              urlProvidedOnStartup = undefined
            } catch (e) {
              console.error(e)
            }
          }
        }
      }

      mainWindow.once('closed', () => {
        mainWindow = undefined
        openMenuItem.label = 'Open Sturdy'
        app.dock?.hide()
      })

      mainWindow.webContents.on('will-navigate', (event, url) => {
        if (!url.startsWith(webURL.href)) {
          shell.openExternal(url)
          event.preventDefault()
        }
      })

      // open target="_blank" links eternally
      mainWindow.webContents.setWindowOpenHandler(({ url }) => {
        shell.openExternal(url)
        return { action: 'deny' }
      })

      // If the page fails to load, display app-fail.html
      mainWindow.webContents.on(
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
          console.error('did-fail-load', { errorCode, errorDescription })
          mainWindow?.loadFile(resourcePath('app-fail.html'), {
            query: {
              goto: webURL.href,
            },
          })
        }
      )

      mutagenManager.setMainWindow(mainWindow)
      longLivedInstances.add(mainWindow)

      try {
        await loadURLWithoutThrowingOnRedirects(mainWindow, url.href)
      } catch (e) {
        console.error('failed to loadURL', e)
      }
    }

    app.on('open-url', async (event, url) => {
      try {
        console.log('open-url', url)
        if (!url.startsWith(protocol + '://')) {
          return
        }
        event.preventDefault()

        const sturdyUrl = new URL(url)
        const newUrl = new URL(sturdyUrl.pathname + sturdyUrl.search, webURL)

        await triggerWindow(newUrl)
      } catch (e) {
        console.error(e)
      }
    })

    app.on('second-instance', (event, commandLine, workingDirectory) => {
      console.log('second-instance', commandLine)

      // Windows handling for opening protocol links while there is already a window open
      if (process.platform === 'win32') {
        const argWithUrl = commandLine.find((arg) => arg.indexOf(protocol) > -1)
        if (argWithUrl) {
          const sturdyUrl = new URL(argWithUrl)

          const newUrl = new URL(sturdyUrl.pathname + sturdyUrl.search, webURL)

          triggerWindow(newUrl)
        } else {
          triggerWindow()
        }
        return
      }

      // Not windows
      triggerWindow()
    })

    app.on('activate', () => {
      // On macOS it's common to re-create a window in the app when the
      // dock icon is clicked and there are no other windows open.
      if (process.platform === 'darwin') {
        triggerWindow()
      }
    })

    await triggerWindow()
  }

  async function quit() {
    try {
      postHog.flush()
      await mutagenManager.cleanup()
      mutagenExecutable.abort()
      mutagenLog.end()
    } finally {
      process.exit()
    }
  }

  main().catch(async (e) => {
    console.error(e)
    try {
      postHog.flush()
      await mutagenManager.cleanup()
      mutagenExecutable.abort()
      mutagenLog.end()
    } catch (er) {
      console.error(er)
    } finally {
      process.exit(1)
    }
  })
}

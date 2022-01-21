import path from 'path'
import { app, crashReporter, Menu, MenuItem, nativeImage, Tray } from 'electron'
import { Updater } from './Updater'
import { resourceBinary, resourcePath } from './resources'
import * as Sentry from '@sentry/electron'
import { Host } from './application'
import { CaptureConsole } from '@sentry/integrations'
import { Application } from './Application'
import { MutagenExecutable } from './mutagen'
import { ApplicationManager } from './ApplicationManager'
import { Logger } from './Logger'

// Start crash reporter before setting up logging
crashReporter.start({
  companyName: 'Sturdy Sweden AB',
  productName: 'Sturdy',
  ignoreSystemCrashHandler: true,
  submitURL:
    'https://o952367.ingest.sentry.io/api/6075838/minidump/?sentry_key=59a9e2de840941b58b49f82b0732e170',
})

// TODO
// Setup logging to file after crash reporter.
// Object.assign(console, log.functions)
// log.transports.file.resolvePath = () => path.join(logsDir, 'main.log')

if (!app.requestSingleInstanceLock()) {
  app.quit()
}

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

const iconSm = nativeImage.createFromPath(resourcePath('AppIconSm.png'))
const iconSmDisconnected = nativeImage.createFromPath(resourcePath('AppIconSmDisconnected.png'))
const iconSmTemplate = nativeImage.createFromPath(resourcePath('AppIconSmTemplate.png'))
const iconSmDisconnectedTemplate = nativeImage.createFromPath(
  resourcePath('AppIconSmDisconnectedTemplate.png')
)

const logger = new Logger()

const iconTray = process.platform === 'darwin' ? iconSmTemplate : iconSm
const iconTrayDisconnected =
  process.platform === 'darwin' ? iconSmDisconnectedTemplate : iconSmDisconnected

const postHogToken =
  process.env.STURDY_POSTHOG_API_KEY ?? 'ZuDRoGX9PgxGAZqY4RF9CCJJLpx14h3szUPzm7XBWSg'

const runAutoUpdater = app.isPackaged && !process.env.STURDY_DISABLE_AUTO_UPDATER

const development = new Host(
  'development',
  'Development',
  new URL('http://localhost:8080'),
  new URL('http://localhost:3000/graphql'),
  new URL('ssh://127.0.0.1:2222')
)
const cloud = new Host(
  'cloud',
  'Cloud',
  new URL('https://getsturdy.com'),
  new URL('https://api.getsturdy.com/graphql'),
  new URL('ssh://sync.getsturdy.com')
)
const selfhosted = new Host(
  'selfhosted',
  'Selfhosted',
  new URL('http://localhost:30080'),
  new URL('http://localhost:30080/api/graphql'),
  new URL('ssh://localhost:30022')
)

const knownHosts =
  process.env.STURDY_DEFAULT_BACKEND === 'development'
    ? [cloud, selfhosted, development]
    : [cloud, selfhosted]

const defaultHost = knownHosts.find((h) => h.id === process.env.STURDY_DEFAULT_BACKEND) ?? cloud

const mutagenExecutable = new MutagenExecutable({
  executablePath: resourceBinary('sturdy-sync'),
  logger: logger,
})

const manager = new ApplicationManager(
  knownHosts,
  mutagenExecutable,
  postHogToken,
  app.isPackaged,
  protocol,
  logger
)

let tray: Tray | undefined

const menu = (application: Application) => {
  const menu = new Menu()
  application.status.appendMenuItem(menu)
  menu.append(new MenuItem({ type: 'separator' }))
  menu.append(
    new MenuItem({
      label: 'Open Sturdy',
      click: () => application.open(),
    })
  )
  manager.appendMenu(menu)
  menu.append(
    new MenuItem({
      label: 'Debug',
      submenu: Menu.buildFromTemplate([
        new MenuItem({
          label: 'Force Restart Syncer',
          click: () => application.forceRestart(),
        }),
      ]),
    })
  )
  menu.append(new MenuItem({ type: 'separator' }))
  if (runAutoUpdater) {
    Updater.start(logger, menu).catch((e) => {
      logger.error('Auto updater error!', e)
    })
  }
  menu.append(
    new MenuItem({
      label: 'Quit Sturdy',
      click: async () => {
        try {
          await manager.cleanup()
        } finally {
          process.exit()
        }
      },
      accelerator: 'CommandOrControl+Q',
    })
  )
  return menu
}

manager.on('status', (state) => {
  if (state === 'online') {
    tray?.setImage(iconTray)
  } else {
    tray?.setImage(iconTrayDisconnected)
  }
})

manager.on('switch', async (application) => {
  tray?.setContextMenu(menu(application))

  app.on('window-all-closed', () => {
    // Don't do anything
    // Keep the app running in the tray
  })

  app.on('open-url', async (event, url) => {
    try {
      logger.log('open-url', url)
      if (!url.startsWith(protocol + '://')) {
        return
      }
      event.preventDefault()

      const sturdyUrl = new URL(url)
      const newUrl = new URL(sturdyUrl.pathname + sturdyUrl.search, application.host.webURL)

      await application.open(newUrl)
    } catch (e) {
      logger.error(e)
    }
  })

  app.on('second-instance', (event, commandLine, workingDirectory) => {
    logger.log('second-instance', commandLine)

    // Windows handling for opening protocol links while there is already a window open
    if (process.platform === 'win32') {
      const argWithUrl = commandLine.find((arg) => arg.indexOf(protocol) > -1)
      if (argWithUrl) {
        const sturdyUrl = new URL(argWithUrl)
        const newUrl = new URL(sturdyUrl.pathname + sturdyUrl.search, application.host.webURL)
        application.open(newUrl)
      } else {
        application.open()
      }
      return
    }

    // Not windows
    application.open()
  })

  app.on('activate', () => {
    // On macOS it's common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (process.platform === 'darwin') {
      application.open()
    }
  })

  await application.open()
})

async function main() {
  if (app.isPackaged) {
    await Updater.finalizePendingUpdate()
  }

  await app.whenReady()
  tray = new Tray(iconTrayDisconnected)

  // Adds support for notifications on Windows
  if (process.platform === 'win32') {
    app.setAppUserModelId('com.getsturdy.sturdy')
  }

  manager.refresh()
  // kick off the app with the default host
  manager.set(defaultHost)
}

main().catch(async (e) => {
  logger.error(e)
  try {
    await manager.cleanup()
    mutagenExecutable.abort()
  } catch (er) {
    logger.error(er)
  } finally {
    process.exit(1)
  }
})

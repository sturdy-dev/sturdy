import path from 'path'
import { app, crashReporter, Menu, MenuItem, nativeImage, Tray } from 'electron'
import { Updater } from './Updater'
import { dataPath, resourcePath } from './resources'
import * as Sentry from '@sentry/electron'
import { Status } from './application'
import { CaptureConsole } from '@sentry/integrations'
import { Application } from './Application'
import { ApplicationManager } from './ApplicationManager'
import { Logger } from './Logger'
import log from 'electron-log'
import { Preferences } from './preferences'

// Start crash reporter before setting up logging
crashReporter.start({
  companyName: 'Sturdy Sweden AB',
  productName: 'Sturdy',
  ignoreSystemCrashHandler: true,
  submitURL:
    'https://o952367.ingest.sentry.io/api/6075838/minidump/?sentry_key=59a9e2de840941b58b49f82b0732e170',
})

const logsDir = dataPath('logs')

// Setup logging to file after crash reporter.
Object.assign(console, log.functions)
log.transports.file.resolvePath = () => path.join(logsDir, 'main.log')

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

const status = new Status(logger)

const manager = new ApplicationManager(
  postHogToken,
  app.isPackaged,
  protocol,
  logger,
  status,
  logsDir
)

let tray: Tray | undefined

let preferences: Preferences | undefined

// todo: assemble menu in a better way, this uses a lot of globals
const menu = (application: Application) => {
  const menu = new Menu()
  status.appendMenuItem(menu)
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
  preferences?.appendMenuItem(menu)
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

status.on('change', (state) => {
  if (state === 'online') {
    tray?.setImage(iconTray)
  } else {
    tray?.setImage(iconTrayDisconnected)
  }
})

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

    await manager.open(url)
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
      manager.open(argWithUrl)
    } else {
      manager.open()
    }
  } else {
    manager.open()
  }
})

app.on('activate', () => {
  // On macOS it's common to re-create a window in the app when the
  // dock icon is clicked and there are no other windows open.
  if (process.platform === 'darwin') {
    manager.open()
  }
})

let application: Application | undefined

manager.on('switch', async (application: Application) => {
  application = application
  tray?.setContextMenu(menu(application))
  await application.open()
})

async function main() {
  if (app.isPackaged) {
    await Updater.finalizePendingUpdate()
  }

  preferences = await Preferences.open(logger)
  await app.whenReady()
  tray = new Tray(iconTrayDisconnected)

  // Adds support for notifications on Windows
  if (process.platform === 'win32') {
    app.setAppUserModelId('com.getsturdy.sturdy')
  }

  preferences.on('open', (host) => manager.set(host))
  manager.on('openPreferences', () => preferences?.showWindow())

  preferences.on('hostsChanged', (hosts) => {
    manager.updateHosts(hosts)
    tray?.setContextMenu(menu(application!))
  })

  await manager.updateHosts(preferences.hosts)

  // kick off the app with the default host
  manager.set(preferences.hosts[0])
}

main().catch(async (e) => {
  logger.error(e)
  try {
    await manager.cleanup()
  } catch (er) {
    logger.error(er)
  } finally {
    process.exit(1)
  }
})

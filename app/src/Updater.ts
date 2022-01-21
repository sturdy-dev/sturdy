import { autoUpdater } from 'electron-updater'
import { MenuItem, Menu } from 'electron'
import { spawn } from 'child_process'
import { app } from 'electron'
import { Logger } from './Logger'

export class Updater {
  readonly #checkForUpdatesMenuItem: MenuItem
  readonly #checkingForUpdatesMenuItem: MenuItem
  readonly #restartToApplyMenuItem: MenuItem
  readonly #logger: Logger

  #interval?: NodeJS.Timer

  private constructor(logger: Logger, menu: Menu) {
    this.#logger = logger.withPrefix('updater')

    this.#checkForUpdatesMenuItem = new MenuItem({
      label: 'Check for Updates',
      click: this.checkForUpdates.bind(this),
    })

    this.#checkingForUpdatesMenuItem = new MenuItem({
      label: 'Checking for Updates...',
      enabled: false,
      visible: false,
    })

    this.#restartToApplyMenuItem = new MenuItem({
      label: 'Restart to Apply Update',
      enabled: false,
      visible: false,
    })

    menu.append(this.#checkForUpdatesMenuItem)
    menu.append(this.#checkingForUpdatesMenuItem)
    menu.append(this.#restartToApplyMenuItem)
  }

  static async start(logger: Logger, menu: Menu): Promise<Updater> {
    const updater = new Updater(logger, menu)
    updater.#listen()
    await updater.checkForUpdates()
    return updater
  }

  static async finalizePendingUpdate() {
    if (process.platform === 'darwin') {
      let shipItPid: number | undefined
      console.log('starting ps process')
      const psProcess = spawn('ps', ['x', '-o', 'pid,command'], {
        stdio: ['ignore', 'pipe', 'ignore'],
      })
      for await (const chunk of psProcess.stdout) {
        const pid = (chunk as Buffer)
          .toString()
          .split('\n')
          .map((s) => s.trim())
          .filter(Boolean)
          .find((line) =>
            line.includes('/Sturdy.app/Contents/Frameworks/Squirrel.framework/Resources/ShipIt')
          )
          ?.split(/ /)?.[0]
        if (pid != null) {
          shipItPid = Number(pid)
          break
        }
      }
      psProcess.kill()

      if (shipItPid != null) {
        while (true) {
          try {
            process.kill(shipItPid, 0)
            console.log('shipIt process lives')
            await new Promise<void>((r) => setTimeout(r, 1000))
          } catch {
            console.log('shipIt process died')
            break
          }
        }

        console.log('relaunching')
        app.relaunch()
        console.log('exiting')
        process.exit()
      } else {
        console.log('no shipIt process found')
      }
    }
  }

  #listen() {
    autoUpdater.on('checking-for-update', () => {
      this.#checkForUpdatesMenuItem.visible = false
      this.#checkingForUpdatesMenuItem.visible = true
      this.#restartToApplyMenuItem.visible = false
    })
    autoUpdater.on('update-not-available', () => {
      this.#checkForUpdatesMenuItem.visible = true
      this.#checkingForUpdatesMenuItem.visible = false
      this.#restartToApplyMenuItem.visible = false
    })
    autoUpdater.on('download-progress', (progress) => {
      // this.#checkForUpdatesMenuItem.visible = true
      // this.#checkingForUpdatesMenuItem.visible = false
      // this.#restartToApplyMenuItem.visible = false
    })
    autoUpdater.on('update-downloaded', () => {
      this.#checkForUpdatesMenuItem.visible = false
      this.#checkingForUpdatesMenuItem.visible = false
      this.#restartToApplyMenuItem.visible = true
    })
    this.#interval = setInterval(this.checkForUpdates.bind(this), 10 * 60_000)
  }

  async checkForUpdates() {
    const result = await autoUpdater.checkForUpdatesAndNotify()
    if (result != null && this.#interval != null) {
      clearInterval(this.#interval)
    }
  }
}

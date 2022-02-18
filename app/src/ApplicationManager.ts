import { Menu, MenuItem } from 'electron'
import { Application, Host, State, Status } from './application'
import { TypedEventEmitter } from './TypedEventEmitter'
import { MutagenDaemon, MutagenExecutable } from './mutagen'
import { Logger } from './Logger'
import { dataPath, resourceBinary } from './resources'
import { createWriteStream } from 'fs'
import path from 'path'
import { WriteStream } from 'fs'

export interface ApplicationManagerEvents {
  switch: [application: Application]
  status: [state: State]
  openPreferences: []
}

export class ApplicationManager extends TypedEventEmitter<ApplicationManagerEvents> {
  readonly #menuItem: MenuItem = new MenuItem({
    label: 'Server',
    submenu: Menu.buildFromTemplate([]),
    visible: false,
  })
  readonly #applications: Map<string, Application> = new Map()

  #activeApplication?: string

  readonly #mutagenExecutable: MutagenExecutable
  readonly #postHogToken: string
  readonly #isAppPackaged: boolean
  readonly #protocol: string
  readonly #logger: Logger
  readonly #daemon: MutagenDaemon
  readonly #mutagenLog: WriteStream
  readonly #status: Status

  constructor(
    postHogToken: string,
    isAppPackaged: boolean,
    protocol: string,
    logger: Logger,
    status: Status,
    logsDir: string
  ) {
    super()

    logger = logger.withPrefix('apps')

    const mutagenLog = createWriteStream(path.join(logsDir, 'mutagen.log'), {
      flags: 'a',
      mode: 0o666,
    })
    logger.log('mutagen log', mutagenLog.path)
    const mutagenDataDirectoryPath = dataPath('mutagen')
    logger.log('mutagen data directory', mutagenDataDirectoryPath)
    const mutagenExecutable = new MutagenExecutable({
      executablePath: resourceBinary('sturdy-sync'),
      logger: logger,
      log: mutagenLog,
      dataDirectory: mutagenDataDirectoryPath,
    })

    const daemon = new MutagenDaemon(logger, mutagenExecutable)
    daemon.on('session-manager-initialized', status.mutagenStarted.bind(status))
    daemon.on('failed-to-start', status.mutagenFailedToStart.bind(status))
    daemon.on('is-running-changed', (isRunning) => {
      if (!isRunning) {
        status.mutagenStopped()
      }
    })
    daemon.on('connection-to-server-dropped', status.connectionDropped.bind(status))

    this.#mutagenExecutable = mutagenExecutable
    this.#postHogToken = postHogToken
    this.#isAppPackaged = isAppPackaged
    this.#protocol = protocol
    this.#logger = logger
    this.#daemon = daemon
    this.#status = status
    this.#mutagenLog = mutagenLog
  }

  #addHosts(hosts: Host[]) {
    this.#menuItem.visible = true
    const submenu = this.#menuItem.submenu!
    hosts.forEach((host) => {
      if (submenu.getMenuItemById(host.id)) return
      const item = new MenuItem({
        id: host.id,
        label: host.title,
        enabled: false,
        click: () => this.set(host),
        type: 'radio',
      })
      submenu.append(item)
      this.#updateHostStatus(host)
    })
  }

  async getOrCreateApplication(host: Host) {
    if (this.#applications.has(host.id)) {
      return this.#applications.get(host.id)!
    }

    const application = await Application.start({
      host,
      mutagenExecutable: this.#mutagenExecutable,
      postHogToken: this.#postHogToken,
      isAppPackaged: this.#isAppPackaged,
      protocol: this.#protocol,
      status: this.#status,
      logger: this.#logger,
      daemon: this.#daemon,
    })
    this.#activeApplication = host.id
    this.#applications.set(host.id, application)
    application.status.on('change', (state) => {
      this.#stateChanged(host, state)
    })
    application.on('openPreferences', () => this.emit('openPreferences'))
    return application
  }

  #stateChanged(host: Host, state: State) {
    if (host.id === this.#activeApplication) {
      this.emit('status', state)
    }
  }

  async set(host: Host) {
    const application = await this.getOrCreateApplication(host)
    const item = this.#menuItem.submenu?.getMenuItemById(host.id)
    if (item) item.checked = true
    this.emit('switch', application)

    // close all other applications
    for (let [id, app] of this.#applications) {
      if (id === host.id) {
        continue
      }
      app.close()
    }
  }

  async open(url?: string) {
    if (!this.#activeApplication) {
      return
    }
    const application = this.#applications.get(this.#activeApplication)!
    if (url) {
      const sturdyUrl = new URL(url)
      const newUrl = new URL(sturdyUrl.pathname + sturdyUrl.search, application.host.webURL)
      await application.open(newUrl)
    } else {
      await application.open()
    }
  }

  appendMenu(menu: Menu) {
    menu.append(this.#menuItem)
  }

  async updateHosts(hosts: Host[]) {
    this.#addHosts(hosts)
    // fetch status of all known hosts
    const promises = []
    for (const host of hosts) {
      promises.push(this.#updateHostStatus(host))
    }
    await Promise.all(promises)
  }

  async #updateHostStatus(host: Host) {
    await host.isUp().then((isUp) => {
      const submenu = this.#menuItem.submenu
      if (!submenu) return
      const item = submenu.getMenuItemById(host.id)
      if (!item) return
      item.enabled = isUp
    })
  }

  async cleanup() {
    this.#logger.log('cleaning up...')
    await Promise.all(Array.from(this.#applications.values()).map((ctx) => ctx.cleanup()))
    this.#mutagenExecutable.abort()
    this.#mutagenLog.end()
  }
}

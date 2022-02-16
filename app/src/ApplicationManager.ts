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
  #hosts: Map<string, Host> = new Map()
  #menuItems: Map<string, MenuItem> = new Map()
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

  #addHost(host: Host) {
    if (this.#hosts.has(host.id)) {
      return
    }
    this.#hosts.set(host.id, host)
    this.#menuItems.set(
      host.id,
      new MenuItem({
        label: host.title,
        enabled: false,
        type: 'radio',
        click: () => this.set(host),
      })
    )
  }

  #addHosts(hosts: Host[]) {
    hosts.forEach((host) => this.#addHost(host))
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
    this.#menuItems.forEach((item) => (item.checked = false))
    this.#menuItems.get(host.id)!.checked = true
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
    const enabledMenues = Array.from(this.#menuItems.values()).filter((item) => item.enabled)
    if (enabledMenues.length <= 1) {
      // don't show the menu if there is only one backend available
      return
    }
    menu.append(
      new MenuItem({
        label: 'Server',
        submenu: Menu.buildFromTemplate(
          Array.from(this.#menuItems.values()).sort((a, b) => a.label.localeCompare(b.label))
        ),
      })
    )
  }

  async updateHosts(hosts: Host[]) {
    this.#addHosts(hosts)
    // fetch status of all known hosts
    const promises = []
    for (const host of hosts) {
      promises.push(
        host.isUp().then((isUp) => {
          this.#menuItems.get(host.id)!.enabled = isUp
        })
      )
    }
    await Promise.all(promises)
  }

  async cleanup() {
    this.#logger.log('cleaning up...')
    await Promise.all(Array.from(this.#applications.values()).map((ctx) => ctx.cleanup()))
    this.#mutagenExecutable.abort()
    this.#mutagenLog.end()
  }
}

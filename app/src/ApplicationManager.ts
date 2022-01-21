import { Menu, MenuItem } from 'electron'
import { Application, Host, State } from './application'
import { TypedEventEmitter } from './TypedEventEmitter'
import { MutagenExecutable } from './mutagen'
import { Logger } from './Logger'

export interface ApplicationManagerEvents {
  switch: [application: Application]
  status: [state: State]
}

export class ApplicationManager extends TypedEventEmitter<ApplicationManagerEvents> {
  readonly #hosts: Map<string, Host>
  readonly #menuItems: Map<string, MenuItem>
  readonly #applications: Map<string, Application> = new Map()

  #activeApplication?: string

  readonly #mutagenExecutable: MutagenExecutable
  readonly #postHogToken: string
  readonly #isAppPackaged: boolean
  readonly #protocol: string
  readonly #logger: Logger

  constructor(
    hosts: Host[],
    mutagenExecutable: MutagenExecutable,
    postHogToken: string,
    isAppPackaged: boolean,
    protocol: string,
    logger: Logger
  ) {
    super()

    this.#hosts = new Map(hosts.map((host) => [host.id, host]))

    this.#menuItems = new Map(
      hosts.map((host) => [
        host.id,
        new MenuItem({
          label: host.title,
          enabled: false,
          type: 'radio',
          click: () => this.set(host),
        }),
      ])
    )

    this.#mutagenExecutable = mutagenExecutable
    this.#postHogToken = postHogToken
    this.#isAppPackaged = isAppPackaged
    this.#protocol = protocol
    this.#logger = logger.withPrefix('apps')
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
      logger: this.#logger,
    })
    this.#activeApplication = host.id
    this.#applications.set(host.id, application)
    application.status.on('change', (state) => {
      this.#stateChanged(host, state)
    })
    return application
  }

  #stateChanged(host: Host, state: State) {
    if (host.id === this.#activeApplication) {
      this.emit('status', state)
    }
  }

  async set(host: Host) {
    const application = await this.getOrCreateApplication(host)
    this.#menuItems.get(host.id)!.checked = true
    this.emit('switch', application)
    for (let [id, app] of this.#applications) {
      if (id === host.id) {
        continue
      }
      app.close()
    }
  }

  appendMenu(menu: Menu) {
    menu.append(
      new MenuItem({
        label: 'Server',
        submenu: Menu.buildFromTemplate(
          Array.from(this.#menuItems.values()).sort((a, b) => a.label.localeCompare(b.label))
        ),
      })
    )
  }

  refresh() {
    this.#hosts.forEach((host) =>
      host.isUp().then((isUp) => {
        this.#menuItems.get(host.id)!.enabled = isUp
      })
    )
  }

  async cleanup() {
    this.#logger.log('cleaning up...')
    await Promise.all(Array.from(this.#applications.values()).map((ctx) => ctx.cleanup()))
  }
}

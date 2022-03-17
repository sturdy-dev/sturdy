import { app, BrowserWindow, ipcMain, Menu, MenuItem } from 'electron'
import { Logger } from '../Logger'
import path from 'path'
import { readFile, writeFile } from 'fs/promises'
import { dataPath } from '../resources'
import { Host } from '../application'
import { TypedEventEmitter } from '../TypedEventEmitter'
import ElectronWindowState from 'electron-window-state'

type DetailedHostConfig = {
  webURL: string
  apiURL: string
  syncURL: string
}

type ShortHostConfig = {
  host: string
}

export type HostConfig = (DetailedHostConfig | ShortHostConfig) & {
  title: string
}

const validateHostConfig = (hostConfig: HostConfig): HostConfig => {
  if (!hostConfig.title) {
    throw new Error(`Title can not be empty`)
  }
  const shortConfig = hostConfig as ShortHostConfig
  if (shortConfig.host) {
    validateShortHostConfig(shortConfig)
  } else {
    validateDetailedHostConfig(hostConfig as DetailedHostConfig)
  }
  return hostConfig
}

const validateShortHostConfig = (hostConfig: ShortHostConfig): ShortHostConfig => {
  if (hostConfig.host.startsWith('http')) {
    try {
      new URL(hostConfig.host)
    } catch (e) {
      throw new Error(`${hostConfig.host} is not a valid host`)
    }
  }
  try {
    new URL(`http://${hostConfig.host}`)
  } catch (e) {
    throw new Error(`${hostConfig.host} is not a valid host`)
  }
  return hostConfig
}

const isValidShortHostConfig = (hostConfig: ShortHostConfig): boolean => {
  try {
    validateShortHostConfig(hostConfig)
    return true
  } catch {
    return false
  }
}

const validateDetailedHostConfig = (hostConfig: DetailedHostConfig): DetailedHostConfig => {
  try {
    new URL(hostConfig.webURL)
  } catch (e) {
    throw new Error(`${hostConfig.webURL} is not a valid host`)
  }
  try {
    new URL(hostConfig.apiURL)
  } catch (e) {
    throw new Error(`${hostConfig.apiURL} is not a valid host`)
  }
  try {
    new URL(hostConfig.syncURL)
  } catch (e) {
    throw new Error(`${hostConfig.syncURL} is not a valid host`)
  }
  return hostConfig
}

type migration = (cfg: HostConfig[]) => HostConfig[]

const updateSelfHosted: migration = (cfg: HostConfig[]) =>
  cfg.map((host) => {
    if (host.title !== 'Self Hosted') {
      return host
    }
    return { title: host.title, host: 'locahost:30080' }
  })

const migrations: Record<string, migration> = {
  'update-default-selfhosted-url': updateSelfHosted,
}

type Config = {
  hosts: HostConfig[]
  appliedMigrations?: string[]
}

const development: HostConfig = {
  title: 'Development',
  webURL: 'http://localhost:8080',
  apiURL: 'http://localhost:3000',
  syncURL: 'ssh://localhost:2222',
}

const cloud: HostConfig = {
  title: 'Cloud',
  webURL: 'https://getsturdy.com',
  apiURL: 'https://api.getsturdy.com',
  syncURL: 'ssh://sync.getsturdy.com',
}

const selfhosted: HostConfig = {
  title: 'Self-hosted',
  host: 'localhost:30080',
}

const httpsPort = '443',
  httpPort = '80'

const hostFromConfig = (hostConfig: HostConfig): Host => {
  if (isValidShortHostConfig(hostConfig as ShortHostConfig)) {
    const config = hostConfig as ShortHostConfig

    let webApiProto = 'http:'

    let baseHost = config.host
    if (config.host.startsWith('https://') || config.host.startsWith('http://')) {
      const webURL = new URL(`${config.host}`)
      webApiProto = webURL.protocol
    } else {
      baseHost = 'http://' + baseHost
    }

    const configURL = new URL(baseHost)

    const port = configURL.port
      ? configURL.port
      : config.host.startsWith('https://')
      ? httpsPort
      : httpPort

    const webURLStr = `${webApiProto}//${configURL.hostname}:${port}`
    const webURL = new URL(webURLStr)

    const apiURLStr = `${webApiProto}//${configURL.hostname}:${port}/api`
    const apiURL = new URL(apiURLStr)

    const syncURLStr = `ssh://${apiURL.hostname}:${port}`
    const syncURL = new URL(syncURLStr)

    return new Host({
      title: hostConfig.title,
      webURL,
      apiURL,
      syncURL,
    })
  }

  // is detailed config

  const detailedConfig = hostConfig as DetailedHostConfig
  validateDetailedHostConfig(detailedConfig)

  return new Host({
    title: hostConfig.title,
    webURL: new URL(detailedConfig.webURL),
    apiURL: new URL(detailedConfig.apiURL),
    syncURL: new URL(detailedConfig.syncURL),
  })
}

const defaultConfig = {
  hosts:
    process.env.STURDY_DEFAULT_BACKEND === 'development'
      ? [development, cloud, selfhosted]
      : [cloud, selfhosted],
}

const preferencesPath = dataPath('preferences.json')

export interface PreferencesEvents {
  hostsChanged: [hosts: Host[]]
  open: [host: Host]
}

export class Preferences extends TypedEventEmitter<PreferencesEvents> {
  #window?: BrowserWindow
  readonly #logger: Logger
  #config: Config
  #hosts: Host[]

  private constructor(logger: Logger, config: Config) {
    super()

    this.#logger = logger.withPrefix('settings')
    this.#config = config
    this.#hosts = config.hosts.map(hostFromConfig)
  }

  get hosts() {
    return this.#hosts
  }

  static async open(logger: Logger) {
    const cfg = await this.#read(logger).then((cfg) =>
      this.migrate(logger.withPrefix('migrate'), cfg)
    )
    return new Preferences(logger, cfg)
  }

  static async migrate(logger: Logger, cfg: Config): Promise<Config> {
    const migrationsToApply = Object.entries(migrations).filter(
      ([key, _]) => !cfg.appliedMigrations?.includes(key)
    )
    if (migrationsToApply.length === 0) {
      return cfg
    }

    cfg = migrationsToApply.reduce((cfg, [migrationName, migrate]) => {
      logger.log(`applying ${migrationName}`)
      cfg.hosts = migrate(cfg.hosts)
      if (cfg.appliedMigrations) {
        cfg.appliedMigrations.push(migrationName)
      } else {
        cfg.appliedMigrations = [migrationName]
      }
      return cfg
    }, cfg)

    await writeFile(preferencesPath, JSON.stringify(cfg, null, 2))
    logger.log(`write config to ${preferencesPath}`, cfg)
    return cfg
  }

  static async #read(logger: Logger): Promise<Config> {
    try {
      const cfg = JSON.parse(await readFile(preferencesPath, 'utf8'))
      logger.log(`read config from ${preferencesPath}`, cfg)
      return cfg
    } catch (e) {
      if ((e as any).code === 'ENOENT') {
        return defaultConfig
      }
      throw e
    }
  }

  appendMenuItem(menu: Menu) {
    menu.append(
      new MenuItem({
        label: 'Preferences',
        accelerator: 'CommandOrControl+,',
        click: () => {
          this.showWindow()
        },
      })
    )
  }

  showWindow() {
    this.#logger.log('showWindow')
    if (this.#window) {
      this.#window.show()
    } else {
      this.#window = this.#newWindow()
      this.#window.on('closed', () => {
        this.#window = undefined
      })
    }
  }

  async #openHost(hostConfig: HostConfig) {
    const host = hostFromConfig(hostConfig)
    this.emit('open', host)
    this.emit('hostsChanged', this.#hosts)
  }

  async #isHostUp(hostConfig: HostConfig) {
    try {
      validateHostConfig(hostConfig)
    } catch {
      return false
    }
    const host = hostFromConfig(hostConfig)
    this.emit('hostsChanged', this.#hosts)
    return host.isUp()
  }

  async #handleDeleteHostConfig(hostConfig: HostConfig) {
    const index = this.#config.hosts.findIndex((h) => h.title === hostConfig.title)
    if (index === -1) {
      throw new Error('Host config not found')
    }
    this.#config.hosts.splice(index, 1)
    this.#hosts.splice(index, 1)
    await this.#saveConfig(this.#config)
  }

  async #handleAddHostConfig(hostConfig: HostConfig) {
    try {
      validateHostConfig(hostConfig)
    } catch (e) {
      throw e
    }
    this.#config.hosts.push(hostConfig)
    this.#hosts.push(hostFromConfig(hostConfig))
    await this.#saveConfig(this.#config)
  }

  async #saveConfig(config: Config) {
    this.#logger.log(`write config to ${preferencesPath}`, config)
    await writeFile(preferencesPath, JSON.stringify(config, null, 2))
    this.emit('hostsChanged', this.#hosts)
  }

  async updateHostConfigs(hostConfigs: HostConfig[]) {
    hostConfigs.forEach(validateHostConfig)
    this.#config.hosts = hostConfigs
    this.#hosts = hostConfigs.map(hostFromConfig)
    await this.#saveConfig(this.#config)
  }

  get config() {
    return this.#config!
  }

  #newWindow() {
    const windowState = ElectronWindowState({
      defaultHeight: 320,
      defaultWidth: 640,
      file: 'preference-window.json',
    })
    const window = new BrowserWindow({
      minWidth: 640,
      minHeight: 320,
      ...windowState,
      webPreferences: {
        nodeIntegration: true,
        devTools: !app.isPackaged,
        preload: path.join(__dirname, 'preferences/preload.js'),
      },
      frame: false,
      titleBarStyle: 'hidden',
      trafficLightPosition: { x: 16, y: 16 },
    })
    ipcMain.handle('config:hosts:list', () => this.#config.hosts)
    ipcMain.handle('config:hosts:add', (_, hostConfig) => this.#handleAddHostConfig(hostConfig))
    ipcMain.handle('config:hosts:delete', (_, hostConfig) =>
      this.#handleDeleteHostConfig(hostConfig)
    )
    ipcMain.handle('config:hosts:isUp', (_, hostConfig) => this.#isHostUp(hostConfig))
    ipcMain.handle('config:hosts:open', (_, hostConfig) => this.#openHost(hostConfig))

    ipcMain.handle('minimize', () => window.minimize())
    ipcMain.handle('maximize', () => window.maximize())
    ipcMain.handle('unmaximize', () => window.unmaximize())
    ipcMain.handle('close', () => window.close())
    ipcMain.handle('isMaximized', () => window.isMaximized())
    ipcMain.handle('isMinimized', () => window.isMinimized())
    ipcMain.handle('isNormal', () => window.isNormal())

    window.on('closed', () => {
      ipcMain.removeHandler('config:hosts:add')
      ipcMain.removeHandler('config:hosts:delete')
      ipcMain.removeHandler('config:hosts:list')
      ipcMain.removeHandler('config:hosts:isUp')
      ipcMain.removeHandler('config:hosts:open')

      ipcMain.removeHandler('minimize')
      ipcMain.removeHandler('maximize')
      ipcMain.removeHandler('unmaximize')
      ipcMain.removeHandler('close')
      ipcMain.removeHandler('isMaximized')
      ipcMain.removeHandler('isMinimized')
      ipcMain.removeHandler('isNormal')
    })
    window.loadFile(path.join(__dirname, 'preferences/index.html'))
    return window
  }
}

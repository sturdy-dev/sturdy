import { app, BrowserWindow, globalShortcut, ipcMain, Menu, MenuItem } from 'electron'
import { Logger } from '../Logger'
import path from 'path'
import { readFile, writeFile } from 'fs/promises'
import { dataPath } from '../resources'
import { Host } from '../application'
import { TypedEventEmitter } from '../TypedEventEmitter'

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
    throw new Error('Host config must have a title')
  }
  try {
    validateShortHostConfig(hostConfig as ShortHostConfig)
  } catch (e) {
    validateDetailedHostConfig(hostConfig as DetailedHostConfig)
  } finally {
    return hostConfig
  }
}

const validateShortHostConfig = (hostConfig: ShortHostConfig): ShortHostConfig => {
  if (!hostConfig.host) {
    throw new Error('Host config must have a host')
  }
  return hostConfig
}

const validateDetailedHostConfig = (hostConfig: DetailedHostConfig): DetailedHostConfig => {
  try {
    new URL(hostConfig.webURL)
  } catch (e) {
    throw new Error('Host config must have a valid webURL')
  }
  try {
    new URL(hostConfig.apiURL)
  } catch (e) {
    throw new Error('Host config must have a valid apiURL')
  }
  try {
    new URL(hostConfig.syncURL)
  } catch (e) {
    throw new Error('Host config must have a valid syncURL')
  }
  return hostConfig
}

type migration = (cfg: Config) => Config

const updateSelfHosted = (cfg: Config): Config => {
  return {
    hosts: [
      ...cfg.hosts.filter((h) => h.title !== 'Self-hosted'),
      {
        title: 'Self-hosted',
        host: 'localhost:30080',
      },
    ],
  }
}

const migrations: migration[] = [updateSelfHosted]

type Config = {
  hosts: HostConfig[]
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

const hostFromConfig = (hostConfig: HostConfig): Host => {
  try {
    const config = hostConfig as ShortHostConfig
    validateShortHostConfig(config)
    return new Host({
      title: hostConfig.title,
      webURL: new URL(`http://${config.host}`),
      apiURL: new URL(`http://${config.host}/api`),
      syncURL: new URL(`ssh://${config.host}`),
    })
  } catch (e) {
    const config = hostConfig as DetailedHostConfig
    validateDetailedHostConfig(config)
    return new Host({
      title: hostConfig.title,
      webURL: new URL(config.webURL),
      apiURL: new URL(config.apiURL),
      syncURL: new URL(config.syncURL),
    })
  }
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

    globalShortcut.register('CommandOrControl+,', () => this.#showWindow())
  }

  get hosts() {
    return this.#hosts
  }

  static async open(logger: Logger) {
    let cfg = await this.#read()
    cfg = await this.migrate(cfg)
    return new Preferences(logger, cfg)
  }

  static async migrate(cfg: Config) {
    for (const migration of migrations) {
      cfg = migration(cfg)
    }
    await writeFile(preferencesPath, JSON.stringify(cfg, null, 2))
    return cfg
  }

  static async #read(): Promise<Config> {
    try {
      return JSON.parse(await readFile(preferencesPath, 'utf8'))
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
        accelerator: process.platform === 'darwin' ? 'Cmd+,' : 'Ctrl+i',
        click: () => this.#showWindow(),
      })
    )
  }

  #showWindow() {
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
    console.log('openHost', hostConfig)
    const host = hostFromConfig(hostConfig)
    console.log('host', host)
    this.emit('open', host)
    this.emit('hostsChanged', this.#hosts)
  }

  async #isHostUp(hostConfig: HostConfig) {
    const host = hostFromConfig(validateHostConfig(hostConfig))
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
    this.emit('hostsChanged', this.#hosts)
    await this.#saveConfig(this.#config)
  }

  async #handleAddHostConfig(hostConfig: HostConfig) {
    hostConfig = validateHostConfig(hostConfig)
    this.#config.hosts.push(hostConfig)
    this.#hosts.push(hostFromConfig(hostConfig))
    this.emit('hostsChanged', this.#hosts)
    await this.#saveConfig(this.#config)
  }

  async #saveConfig(config: Config) {
    this.emit('hostsChanged', this.#hosts)
    await writeFile(preferencesPath, JSON.stringify(config, null, 2))
  }

  get config() {
    return this.#config!
  }

  #newWindow() {
    const window = new BrowserWindow({
      width: 640,
      height: 320,
      minWidth: 640,
      maxWidth: 640,
      minHeight: 230,
      webPreferences: {
        nodeIntegration: true,
        devTools: !app.isPackaged,
        preload: path.join(__dirname, 'preferences/preload.js'),
      },
      titleBarStyle: 'hidden',
      titleBarOverlay: {
        color: '#F9FAFB',
        symbolColor: '#1F2937',
      },
      trafficLightPosition: { x: 16, y: 16 },
    })
    ipcMain.handle('config:hosts:list', () => this.#config.hosts)
    ipcMain.on('config:hosts:add', (_, hostConfig) => this.#handleAddHostConfig(hostConfig))
    ipcMain.on('config:hosts:delete', (_, hostConfig) => this.#handleDeleteHostConfig(hostConfig))
    ipcMain.handle('config:hosts:isUp', (_, hostConfig) => this.#isHostUp(hostConfig))
    ipcMain.handle('config:hosts:open', (_, hostConfig) => this.#openHost(hostConfig))
    window.on('closed', () => {
      ipcMain.removeHandler('config:hosts:add')
      ipcMain.removeHandler('config:hosts:delete')
      ipcMain.removeHandler('config:hosts:list')
      ipcMain.removeHandler('config:hosts:isUp')
      ipcMain.removeHandler('config:hosts:open')
    })
    window.loadFile(path.join(__dirname, 'preferences/index.html'))
    return window
  }
}

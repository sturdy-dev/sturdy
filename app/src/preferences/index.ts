import { app, BrowserWindow, ipcMain, Menu, MenuItem } from 'electron'
import { Logger } from '../Logger'
import path from 'path'
import { readFile, writeFile } from 'fs/promises'
import { dataPath } from '../resources'
import { Host } from '../application'
import { TypedEventEmitter } from '../TypedEventEmitter'

export type HostConfig = {
  title: string
  webURL: string
  apiURL: string
  syncURL: string
  reposBasePath: string
}

const validateHostConfig = (hostConfig: HostConfig): HostConfig => {
  if (!hostConfig.title) {
    throw new Error('Host config must have a title')
  }
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
  if (!hostConfig.reposBasePath) {
    throw new Error('Host config must have a reposBasePath')
  }
  return hostConfig
}

type Config = {
  hosts: HostConfig[]
}

const development: HostConfig = {
  title: 'Development',
  webURL: 'http://localhost:8080',
  apiURL: 'http://localhost:3000',
  syncURL: 'ssh://127.0.0.1:2222',
  reposBasePath: '/repos',
}

const cloud: HostConfig = {
  title: 'Cloud',
  webURL: 'https://getsturdy.com',
  apiURL: 'https://api.getsturdy.com',
  syncURL: 'ssh://sync.getsturdy.com',
  reposBasePath: '/repos',
}
const selfhosted: HostConfig = {
  title: 'Self-hosted',
  webURL: 'http://localhost:30080',
  apiURL: 'http://localhost:30080/api',
  syncURL: 'ssh://localhost:30022',
  reposBasePath: '/var/data/repos',
}

const hostFromConfig = (config: HostConfig): Host => {
  return new Host({
    title: config.title,
    webURL: new URL(config.webURL),
    apiURL: new URL(config.apiURL),
    syncURL: new URL(config.syncURL),
    reposBasePath: config.reposBasePath,
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
    return new Preferences(logger, await this.#read())
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
  }

  async #isHostUp(hostConfig: HostConfig) {
    const host = hostFromConfig(validateHostConfig(hostConfig))
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
    await writeFile(preferencesPath, JSON.stringify(config, null, 2))
  }

  get config() {
    return this.#config!
  }

  #newWindow() {
    const window = new BrowserWindow({
      width: 1174,
      height: 460,
      minWidth: 1174,
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
    ipcMain.handle('config:get', () => this.#config)
    ipcMain.on('config:hosts:add', (_, hostConfig) => this.#handleAddHostConfig(hostConfig))
    ipcMain.on('config:hosts:delete', (_, hostConfig) => this.#handleDeleteHostConfig(hostConfig))
    ipcMain.handle('config:hosts:isUp', (_, hostConfig) => this.#isHostUp(hostConfig))
    ipcMain.handle('config:hosts:open', (_, hostConfig) => this.#openHost(hostConfig))
    window.on('closed', () => {
      ipcMain.removeHandler('config:hosts:add')
      ipcMain.removeHandler('config:hosts:delete')
      ipcMain.removeHandler('config:get')
      ipcMain.removeHandler('config:hosts:isUp')
      ipcMain.removeHandler('config:hosts:open')
    })
    window.loadFile(path.join(__dirname, 'preferences/index.html'))
    return window
  }
}

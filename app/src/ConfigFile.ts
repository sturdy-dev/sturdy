import { readFile, writeFile, rm } from 'fs/promises'
import { Config, ViewConfig } from './Config'
import { join } from 'path'
import { homedir } from 'os'

export class ConfigFile {
  readonly #path: string
  #data!: Config

  private constructor(path: string) {
    this.#path = path
  }

  static async open(path: string) {
    const file = new ConfigFile(path)
    await file.#read()
    return file
  }

  async migrate(): Promise<boolean> {
    const oldPath = process.env.STURDY_OLD_CONFIG_PATH || join(homedir(), '.sturdy')

    let oldConfig: any
    try {
      oldConfig = JSON.parse(await readFile(oldPath, 'utf-8'))
    } catch {
      // No old config found. Nothing to do.
      return false
    }

    if (!oldConfig.views) {
      return false
    }

    await this.update((data) => {
      data.views.push(
        // Read old config to find views
        ...((oldConfig.views as Partial<ViewConfig>[] | undefined)
          // Only keep valid configs...
          ?.filter((v): v is ViewConfig => typeof v.id === 'string' && typeof v.path === 'string')
          // ... that haven't already been migrated
          .filter(
            (v) => !data.views.some((existing) => existing.id === v.id && existing.path === v.path)
          ) ?? [])
      )
    })

    oldConfig['migrated-views'] = oldConfig.views
    delete oldConfig.views

    await writeFile(oldPath, JSON.stringify(oldConfig, null, 4))
    return true
  }

  async #read() {
    try {
      this.#data = JSON.parse(await readFile(this.#path, 'utf-8'))
    } catch (e) {
      if ((e as any).code === 'ENOENT') {
        this.#data = Config.defaultConfig()
        return
      }
      throw e
    }
  }

  get data() {
    return this.#data
  }

  async update(f: (data: Config) => void) {
    let newValue = JSON.parse(JSON.stringify(this.#data))
    f(newValue)
    await writeFile(this.#path, JSON.stringify(newValue, null, 4))
    this.#data = newValue
  }
}

import { Logger } from '../Logger'

export class Host {
  readonly #id: string
  readonly #title: string
  readonly #webURL: URL
  readonly #apiURL: URL
  readonly #syncURL: URL
  readonly #logger: Logger

  constructor(id: string, title: string, webURL: URL, apiURL: URL, syncURL: URL) {
    this.#id = id
    this.#title = title
    this.#webURL = webURL
    this.#apiURL = apiURL
    this.#syncURL = syncURL
    this.#logger = new Logger('host', id)
  }

  get id() {
    return this.#id
  }

  get title() {
    return this.#title
  }

  get webURL() {
    return this.#webURL
  }

  get apiURL() {
    return this.#apiURL
  }

  get syncURL() {
    return this.#syncURL
  }

  async isUp(): Promise<boolean> {
    const fetch = (await import('node-fetch')).default
    const healthcheckURL = new URL('/readyz', this.#apiURL)
    try {
      this.#logger.log('checking if the host is up...')
      const res = await fetch(healthcheckURL.href)
      const isUp = res.ok && res.status === 200
      this.#logger.log(`host is ${isUp ? 'up' : 'down'}`)
      return isUp
    } catch (e) {
      this.#logger.log(`host is down: ${e}`)
      return false
    }
  }
}

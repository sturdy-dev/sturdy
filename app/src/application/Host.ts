import path from 'path'
import { Logger } from '../Logger'

export class Host {
  readonly #id: string
  readonly #title: string
  readonly #webURL: URL
  readonly #apiURL: URL
  readonly #graphqlURL: URL
  readonly #syncURL: URL
  readonly #reposBasePath: string
  readonly #logger: Logger

  constructor({
    title,
    webURL,
    apiURL,
    syncURL,
    reposBasePath,
  }: {
    title: string
    webURL: URL
    apiURL: URL
    syncURL: URL
    reposBasePath: string
  }) {
    this.#id = title.toLowerCase()
    this.#title = title
    this.#webURL = webURL
    this.#apiURL = apiURL
    this.#graphqlURL = new URL(path.join(apiURL.pathname, 'graphql'), apiURL.origin)
    this.#syncURL = syncURL
    this.#logger = new Logger('host', title)
    this.#reposBasePath = reposBasePath
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

  get graphqlURL() {
    return this.#graphqlURL
  }

  get syncURL() {
    return this.#syncURL
  }

  get reposBasePath() {
    return this.#reposBasePath
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

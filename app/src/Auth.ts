import Electron, { Session } from 'electron'
import { EventEmitter } from 'events'
import { decode } from 'jsonwebtoken'
import { TypedEventEmitter } from './TypedEventEmitter'

export interface AuthEvents {
  'logged-in': [jwt: string]
  'logged-out': []
}

export class Auth extends TypedEventEmitter<AuthEvents> {
  static readonly #AUTH_COOKIE_NAME = 'auth'

  readonly #session: Session
  readonly #apiURL: URL
  #jwt?: string

  private constructor(session: Session, apiURL: URL) {
    super()
    this.#session = session
    this.#apiURL = apiURL
  }

  static async start(
    apiURL: URL,
    session: Session = Electron.session.defaultSession
  ): Promise<Auth> {
    const auth = new Auth(session, apiURL)

    await auth.#start()

    return auth
  }

  async #start() {
    const cookies = await this.#session.cookies.get({
      name: Auth.#AUTH_COOKIE_NAME,
      domain: this.#apiURL.hostname,
    })

    for (const cookie of cookies) {
      try {
        const jwt = decode(cookie.value, { json: true })
        if (jwt?.exp == null) {
          continue
        }
        if (jwt.exp * 1000 > Date.now()) {
          this.#jwt = cookie.value
          break
        }
      } catch {}
    }

    this.#session.cookies.on('changed', (_, cookie, __, removed) => {
      if (cookie.name !== Auth.#AUTH_COOKIE_NAME || cookie.domain !== this.#apiURL.hostname) {
        return
      }

      if (removed) {
        this.#jwt = undefined
        this.emit('logged-out')
      } else {
        this.#jwt = cookie.value
        this.emit('logged-in', cookie.value)
      }
    })
  }

  get jwt(): string | undefined {
    return this.#jwt
  }
}

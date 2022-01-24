export class Logger {
  readonly #prefix: string[]

  constructor(...prefix: string[]) {
    this.#prefix = prefix
  }

  withPrefix(...prefix: string[]) {
    return new Logger(...this.#prefix, ...prefix)
  }

  log(...args: any[]) {
    if (this.#prefix.length > 0) {
      console.log(`[${this.#prefix.join('.')}]`, ...args)
    } else {
      console.log(...args)
    }
  }

  error(...args: any[]) {
    if (this.#prefix.length > 0) {
      console.error(`[${this.#prefix.join('.')}]`, ...args)
    } else {
      console.error(...args)
    }
  }
}

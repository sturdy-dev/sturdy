export class Logger {
  readonly #prefix: string[]

  constructor(...prefix: string[]) {
    this.#prefix = prefix
  }

  withPrefix(...prefix: string[]) {
    return new Logger(...this.#prefix, ...prefix)
  }

  log(...args: any[]) {
    console.log(`[${this.#prefix.join('.')}]`, ...args)
  }

  error(...args: any[]) {
    console.error(`[${this.#prefix.join('.')}]`, ...args)
  }
}

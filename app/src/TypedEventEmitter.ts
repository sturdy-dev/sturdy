import EventEmitter from 'events'

export class TypedEventEmitter<T extends { readonly [P in keyof T]: any[] }> {
  #emitter = new EventEmitter()

  setMaxListeners(n: number): this {
    this.#emitter.setMaxListeners(n)
    return this
  }

  protected emit<P extends Exclude<keyof T, number>>(eventName: P, ...args: T[P]): boolean {
    return this.#emitter.emit(eventName, ...args)
  }

  on<P extends Exclude<keyof T, number>>(eventName: P, listener: (...args: T[P]) => void): this {
    this.#emitter.on(eventName, listener as any)
    return this
  }

  once<P extends Exclude<keyof T, number>>(eventName: P, listener: (...args: T[P]) => void): this {
    this.#emitter.once(eventName, listener as any)
    return this
  }

  off<P extends Exclude<keyof T, number>>(eventName: P, listener: (...args: T[P]) => void): this {
    this.#emitter.off(eventName, listener as any)
    return this
  }
}

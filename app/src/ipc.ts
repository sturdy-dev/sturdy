import { State } from './application'

export class IPCMethod<A extends any[], R> {
  readonly #channel: string
  static #CHANNEL_GEN = 0

  constructor(channel: string) {
    this.#channel = channel
  }

  async call(...args: A): Promise<R> {
    const { ipcRenderer } = await import('electron')
    const resolveChannel = String(IPCMethod.#CHANNEL_GEN++)
    const rejectChannel = String(IPCMethod.#CHANNEL_GEN++)
    const promise = new Promise<R>((resolve, reject) => {
      ipcRenderer.once(resolveChannel, (_event, reply: R) => resolve(reply))
      ipcRenderer.once(rejectChannel, (_event, reply: unknown) => reject(reply))
    })
    ipcRenderer.send(this.#channel, resolveChannel, rejectChannel, ...args)
    return promise
  }

  readonly #cleanupFunctions: (() => void)[] = []

  async implement(impl: (...args: A) => R | Promise<R>) {
    const { ipcMain } = await import('electron')
    async function listener(
      event: Electron.IpcMainEvent,
      resolveChannel: string,
      rejectChannel: string,
      ...args: any[]
    ) {
      try {
        const reply = await impl(...(args as A))
        event.reply(resolveChannel, reply)
      } catch (e) {
        event.reply(rejectChannel, e)
      }
    }
    ipcMain.on(this.#channel, listener)
    this.#cleanupFunctions.push(() => ipcMain.off(this.#channel, listener))
  }

  clean() {
    let f: undefined | (() => void)
    while ((f = this.#cleanupFunctions.shift())) {
      f()
    }
  }
}

export var sharedAppIpc = {
  isAuthenticated: new IPCMethod<[], boolean>('isAuthenticated'),
  canGoBack: new IPCMethod<[], boolean>('canGoBack'),
  canGoForward: new IPCMethod<[], boolean>('canGoForward'),
  goBack: new IPCMethod<[], void>('goBack'),
  goForward: new IPCMethod<[], void>('goForward'),
  state: new IPCMethod<[], State>('state'),
  forceRestartMutagen: new IPCMethod<[], void>('forceRestartMutagen'),
  minimize: new IPCMethod<[], void>('minimize'),
  maximize: new IPCMethod<[], void>('maximize'),
  unmaximize: new IPCMethod<[], void>('unmaximize'),
  close: new IPCMethod<[], void>('close'),
  isMaximized: new IPCMethod<[], boolean>('isMaximized'),
  isMinimized: new IPCMethod<[], boolean>('isMinimized'),
  isNormal: new IPCMethod<[], boolean>('isNormal'),
}

// IPC methods that are only available if there is a mutagen session
export var sharedMutagenIpc = {
  createView: new IPCMethod<[workspaceID: string, mountPath: string], string>('createView'),
  createNewViewWithDialog: new IPCMethod<[workspaceID: string, codebaseSlug: string], string>(
    'createNewViewWithDialog'
  ),
}

export type AppIPC = {
  [P in keyof typeof sharedAppIpc]: typeof sharedAppIpc[P] extends IPCMethod<infer A, infer R>
    ? (...args: A) => R | Promise<R>
    : never
}

export type MutagenIPC = {
  [P in keyof typeof sharedMutagenIpc]: typeof sharedMutagenIpc[P] extends IPCMethod<
    infer A,
    infer R
  >
    ? (...args: A) => R | Promise<R>
    : never
}

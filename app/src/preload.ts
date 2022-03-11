import { contextBridge } from 'electron'

import { AppIPC, MutagenIPC, sharedAppIpc, sharedMutagenIpc } from './ipc'

export interface AppEnvironment {
  platform: typeof process.platform
  frameless: boolean
}

const appEnvironment: AppEnvironment = {
  platform: process.platform,
  frameless: true,
}

declare global {
  const mutagenIpc: MutagenIPC
  const ipc: AppIPC
  const appEnvironment: AppEnvironment
}

contextBridge.exposeInMainWorld(
  'mutagenIpc',
  Object.fromEntries(
    Object.entries(sharedMutagenIpc).map(([channel, method]) => {
      return [channel, method.call.bind(method)]
    })
  )
)

contextBridge.exposeInMainWorld(
  'ipc',
  Object.fromEntries(
    Object.entries(sharedAppIpc).map(([channel, method]) => {
      return [channel, method.call.bind(method)]
    })
  )
)

contextBridge.exposeInMainWorld('appEnvironment', appEnvironment)

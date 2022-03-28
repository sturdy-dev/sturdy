import { contextBridge, ipcRenderer } from 'electron'
import { HostConfig } from '.'

contextBridge.exposeInMainWorld('ipc', {
  addHostConfig: async (cfg: HostConfig) => ipcRenderer.invoke('config:hosts:add', cfg),
  deleteHostConfig: async (cfg: HostConfig) => ipcRenderer.invoke('config:hosts:delete', cfg),
  isHostUp: async (cfg: HostConfig) => ipcRenderer.invoke('config:hosts:isUp', cfg),
  listHosts: async () => ipcRenderer.invoke('config:hosts:list'),
  openHost: async (cfg: HostConfig) => ipcRenderer.invoke('config:hosts:open', cfg),

  setChannel: async (channel: string) => ipcRenderer.invoke('config:channel:set', channel),
  getChannel: async () => ipcRenderer.invoke('config:channel:get'),

  minimize: async () => ipcRenderer.invoke('minimize'),
  maximize: async () => ipcRenderer.invoke('maximize'),
  unmaximize: async () => ipcRenderer.invoke('unmaximize'),
  close: async () => ipcRenderer.invoke('close'),
  isMaximized: async () => ipcRenderer.invoke('isMaximized'),
  isMinimized: async () => ipcRenderer.invoke('isMinimized'),
  isNormal: async () => ipcRenderer.invoke('isNormal'),
})

const appEnvironment = {
  platform: process.platform,
  frameless: true,
}

contextBridge.exposeInMainWorld('appEnvironment', appEnvironment)

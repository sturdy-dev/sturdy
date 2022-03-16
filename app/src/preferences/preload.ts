import { contextBridge, ipcRenderer } from 'electron'
import { HostConfig } from '.'

contextBridge.exposeInMainWorld('ipc', {
  addHostConfig: async (cfg: HostConfig) => ipcRenderer.invoke('config:hosts:add', cfg),
  deleteHostConfig: async (cfg: HostConfig) => ipcRenderer.invoke('config:hosts:delete', cfg),
  isHostUp: async (cfg: HostConfig) => ipcRenderer.invoke('config:hosts:isUp', cfg),
  listHosts: async () => ipcRenderer.invoke('config:hosts:list'),
  openHost: async (cfg: HostConfig) => ipcRenderer.invoke('config:hosts:open', cfg),
})

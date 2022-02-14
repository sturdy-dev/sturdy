import { contextBridge, ipcRenderer } from 'electron'
import { HostConfig } from '.'

contextBridge.exposeInMainWorld('ipc', {
  addHostConfig: (cfg: HostConfig) => ipcRenderer.send('config:hosts:add', cfg),
  deleteHostConfig: (cfg: HostConfig) => ipcRenderer.send('config:hosts:delete', cfg),
  isHostUp: async (cfg: HostConfig) => ipcRenderer.invoke('config:hosts:isUp', cfg),
  getConfig: () => ipcRenderer.invoke('config:get'),
  openHost: (cfg: HostConfig) => ipcRenderer.invoke('config:hosts:open', cfg),
})

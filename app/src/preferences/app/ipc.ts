import { Config, HostConfig } from './shims/ipc'

const getConfig = (): Promise<Config> => {
  return window.ipc.getConfig()
}

const addHostConfig = async (cfg: HostConfig): Promise<void> => {
  window.ipc.addHostConfig(cfg)
}

const isHostUp = async (cfg: HostConfig): Promise<boolean> => {
  return window.ipc.isHostUp(cfg)
}

const openHost = async (cfg: HostConfig): Promise<void> => {
  window.ipc.openHost(cfg)
}

const deleteHost = async (cfg: HostConfig): Promise<void> => {
  window.ipc.deleteHostConfig(cfg)
}

export default {
  getConfig,
  addHostConfig,
  isHostUp,
  openHost,
  deleteHost,
}

import { HostConfig } from './shims/ipc'

const listHosts = (): Promise<HostConfig[]> => {
  return window.ipc.listHosts()
}

const addHostConfig = async (cfg: HostConfig): Promise<void> => {
  return window.ipc.addHostConfig(cfg)
}

const isHostUp = async (cfg: HostConfig): Promise<boolean> => {
  return window.ipc.isHostUp(cfg)
}

const openHost = async (cfg: HostConfig): Promise<void> => {
  return window.ipc.openHost(cfg)
}

const deleteHost = async (cfg: HostConfig): Promise<void> => {
  return window.ipc.deleteHostConfig(cfg)
}

const setChannel = async (channel: string): Promise<void> => {
  return window.ipc.setChannel(channel)
}

const getChannel = async (): Promise<string> => {
  return window.ipc.getChannel()
}

export default {
  listHosts,
  addHostConfig,
  isHostUp,
  openHost,
  deleteHost,
  setChannel,
  getChannel,
}

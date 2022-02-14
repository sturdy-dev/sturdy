import { atom } from 'nanostores'

import type { HostConfig } from '../shims/ipc'

export const servers = atom<HostConfig[]>([])

export const add = (server: HostConfig) => {
  servers.set([...servers.get(), server])
}

export const remove = (server: HostConfig) => {
  servers.set(servers.get().filter((s) => s.title !== server.title))
}

export const list = () => {
  return servers.get()
}

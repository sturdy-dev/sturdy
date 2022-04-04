import { atom } from 'nanostores'

import type { HostConfig as cfg } from '../shims/ipc'

export type HostConfig = cfg & {
  isUp: boolean
}

export const servers = atom<HostConfig[]>([])

export const update = (server: HostConfig) => {
  servers.set(servers.get().map((s) => (s.title === server.title ? server : s)))
}

export const add = (server: HostConfig) => {
  servers.set([...servers.get(), server])
}

export const set = (hosts: HostConfig[]) => {
  servers.set(hosts)
}

export const remove = (server: HostConfig) => {
  servers.set(servers.get().filter((s) => s.title !== server.title))
}

export const list = () => {
  return servers.get()
}

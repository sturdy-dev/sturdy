export type HostConfig = (DetailedHostConfig | ShortHostConfig) & {
  title: string
}

export type DetailedHostConfig = {
  webURL: string
  apiURL: string
  syncURL: string
}

export type ShortHostConfig = {
  host: string
}

declare global {
  interface Window {
    readonly ipc: {
      openHost: (cfg: HostConfig) => void
      listHosts: () => Promise<HostConfig[]>
      addHostConfig: (cfg: HostConfig) => void
      deleteHostConfig: (cfg: HostConfig) => void
      isHostUp: (cfg: HostConfig) => Promise<boolean>
    }
  }
}

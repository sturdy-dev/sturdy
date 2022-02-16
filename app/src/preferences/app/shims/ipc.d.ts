export type HostConfig = {
  title: string
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

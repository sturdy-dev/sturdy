export type HostConfig = {
  title: string
  webURL: string
  apiURL: string
  syncURL: string
  reposBasePath: string
}

export type Config = {
  hosts: HostConfig[]
}

declare global {
  interface Window {
    readonly ipc: {
      openHost: (cfg: HostConfig) => void
      getConfig: () => Promise<Config>
      addHostConfig: (cfg: HostConfig) => void
      deleteHostConfig: (cfg: HostConfig) => void
      isHostUp: (cfg: HostConfig) => Promise<boolean>
    }
  }
}

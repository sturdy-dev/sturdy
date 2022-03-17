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
      openHost: (cfg: HostConfig) => Promise<void>
      listHosts: () => Promise<HostConfig[]>
      addHostConfig: (cfg: HostConfig) => Promise<void>
      deleteHostConfig: (cfg: HostConfig) => Promise<void>
      isHostUp: (cfg: HostConfig) => Promise<boolean>

      minimize: () => Promise<void>
      maximize: () => Promise<void>
      unmaximize: () => Promise<void>
      close: () => Promise<void>
      isMinimized: () => Promise<boolean>
      isMaximized: () => Promise<boolean>
      isNormal: () => Promise<boolean>
    }
    readonly appEnvitonment: {
      platform: 'linux' | 'darwin' | 'win32'
      frameless: boolean
    }
  }
}

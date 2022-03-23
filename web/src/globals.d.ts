declare global {
  interface Window {
    ipc?: IPC
    appEnvironment?: AppEnvironment
    mutagenIPC?: MutagenIPC
  }
}

interface IPC {
  isAuthenticated: () => Promise<boolean>
  canGoBack: () => Promise<boolean>
  canGoForward: () => Promise<boolean>
  goBack: () => void
  goForward: () => void
  state: () => 'offline' | 'starting' | 'creating-ssh-key' | 'uploading-ssh-key' | 'online'
  forceRestartMutagen: () => void
  minimize: () => void
  maximize: () => void
  unmaximize: () => void
  close: () => void
  isMinimized: () => Promise<boolean>
  isMaximized: () => Promise<boolean>
  isNormal: () => Promise<boolean>
}

interface AppEnvironment {
  frameless: boolean
  platform: 'linux' | 'darwin' | 'win32'
}

interface MutagenIPC {
  createView: (workspaceId: string, mountPath: string) => Promise<string>
  createNewViewWithDialog: (workspaceId: string) => Promise<string>
  version?: string
}

import { Menu, MenuItem, nativeImage } from 'electron'
import { MutagenSession } from './MutagenSession'
import { resourcePath } from './resources'
import { TypedEventEmitter } from './TypedEventEmitter'

export type AppStatusState =
  | 'offline'
  | 'starting'
  | 'creating-ssh-key'
  | 'uploading-ssh-key'
  | 'online'

const disconnectedIcon = nativeImage.createFromPath(resourcePath('TrayStatusDisconnected.png'))
const connectingIcon = nativeImage.createFromPath(resourcePath('TrayStatusConnecting.png'))
const connectedIcon = nativeImage.createFromPath(resourcePath('TrayStatusConnected.png'))

export interface AppStatusEvents {
  change: [state: AppStatusState]
}

export class AppStatus extends TypedEventEmitter<AppStatusEvents> {
  #state: AppStatusState = 'starting'

  readonly #connectedMenuItem = new MenuItem({
    label: 'Connected',
    enabled: false,
    visible: false,
    icon: connectedIcon,
  })
  readonly #connectingMenuItem = new MenuItem({
    label: 'Connecting',
    enabled: false,
    visible: false,
    icon: connectingIcon,
  })
  readonly #disconnectedMenuItem = new MenuItem({
    label: 'Disconnected',
    enabled: false,
    visible: false,
    icon: disconnectedIcon,
  })

  appendMenuItem(menu: Menu) {
    menu.append(this.#connectedMenuItem)
    menu.append(this.#connectingMenuItem)
    menu.append(this.#disconnectedMenuItem)
  }

  #setState(value: AppStatusState) {
    this.#state = value
    this.emit('change', value)
    console.log('AppStatus', { state: value })

    switch (value) {
      case 'offline':
        this.#disconnectedMenuItem.visible = true
        this.#connectingMenuItem.visible = false
        this.#connectedMenuItem.visible = false
        break
      case 'starting':
      case 'creating-ssh-key':
      case 'uploading-ssh-key':
        this.#disconnectedMenuItem.visible = false
        this.#connectingMenuItem.visible = true
        this.#connectedMenuItem.visible = false
        break
      case 'online':
        this.#disconnectedMenuItem.visible = false
        this.#connectingMenuItem.visible = false
        this.#connectedMenuItem.visible = true
        break
    }
  }

  get state(): AppStatusState {
    return this.#state
  }

  sessionsChanged(sessions: MutagenSession[]) {
    if (sessions.length === 0) {
      // Don't set the status.
      // The sessions list is empty while the sessions are being created (one by one)
      return
    }

    const someSessionHaveNotStarted = sessions.some((s) => !s.isRunning)
    if (someSessionHaveNotStarted) {
      this.#setState('starting')
    } else {
      this.#setState('online')
    }
  }

  mutagenStarted() {
    this.#setState('starting')
  }

  reconciledSessions() {
    this.#setState('online')
  }

  mutagenFailedToStart(_error: Error) {
    this.#setState('offline')
  }

  willCreateSSHKeys() {
    this.#setState('creating-ssh-key')
  }

  willUploadSSHKeys() {
    this.#setState('uploading-ssh-key')
  }

  createdAndUploadedSSHKeys() {
    this.#setState('starting')
  }

  mutagenStopped() {
    this.#setState('offline')
  }

  connectionDropped() {
    this.#setState('offline')
  }
}

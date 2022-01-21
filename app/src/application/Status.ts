import { Menu, MenuItem, nativeImage } from 'electron'
import { Logger } from '../Logger'
import { MutagenSession } from '../mutagen'
import { resourcePath } from '../resources'
import { TypedEventEmitter } from '../TypedEventEmitter'

export type State = 'offline' | 'starting' | 'creating-ssh-key' | 'uploading-ssh-key' | 'online'

const disconnectedIcon = nativeImage.createFromPath(resourcePath('TrayStatusDisconnected.png'))
const connectingIcon = nativeImage.createFromPath(resourcePath('TrayStatusConnecting.png'))
const connectedIcon = nativeImage.createFromPath(resourcePath('TrayStatusConnected.png'))

export interface StatusEvents {
  change: [state: State]
}

export class Status extends TypedEventEmitter<StatusEvents> {
  #state: State = 'starting'
  readonly #logger: Logger

  constructor(logger: Logger) {
    super()
    this.#logger = logger.withPrefix('status')
  }

  readonly #connectedMenuItem = new MenuItem({
    label: 'Connected',
    enabled: true,
    visible: false,
    icon: connectedIcon,
  })
  readonly #connectingMenuItem = new MenuItem({
    label: 'Connecting',
    enabled: true,
    visible: false,
    icon: connectingIcon,
  })
  readonly #disconnectedMenuItem = new MenuItem({
    label: 'Disconnected',
    enabled: true,
    visible: false,
    icon: disconnectedIcon,
  })

  appendMenuItem(menu: Menu) {
    menu.append(this.#connectedMenuItem)
    menu.append(this.#connectingMenuItem)
    menu.append(this.#disconnectedMenuItem)
  }

  #setState(value: State) {
    this.#state = value
    this.emit('change', value)
    this.#logger.log(value)

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

  get state(): State {
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

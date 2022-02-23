import { ChildProcess } from 'child_process'
import { MutagenExecutable } from './Executable'
import { rm } from 'fs/promises'
import { Readable } from 'stream'
import { TypedEventEmitter } from '../TypedEventEmitter'
import { MutagenSessionState } from './Session'
import { Logger } from '../Logger'

interface MutagenDaemonEvents {
  'session-manager-initialized': []
  'failed-to-start': [error: Error]
  'is-running-changed': [isRunning: boolean]
  'session-state-changed': [
    sessionName: string,
    fromState: MutagenSessionState,
    toState: MutagenSessionState
  ]
  'connection-to-server-dropped': []
}

export class MutagenDaemon extends TypedEventEmitter<MutagenDaemonEvents> {
  static readonly #MUTAGEN_STATE_TRANSITION_LOG_PATTERN =
    /([a-z0-9-]+): ([A-Za-z]+) -> ([A-Za-z]+):/

  readonly #executable: MutagenExecutable
  readonly #logger: Logger

  #process?: ChildProcess

  constructor(logger: Logger, executable: MutagenExecutable) {
    super()
    this.#logger = logger.withPrefix('mutagen-daemon')
    this.#executable = executable

    this.setMaxListeners(50)
  }

  get isRunning() {
    return this.#process != null && !this.#process.killed
  }

  #onStdioChunk(chunk: Buffer) {
    const lines = chunk
      .toString()
      .split('\n')
      .map((s) => s.trim())
      .filter(Boolean) // Remove empty lines

    for (const line of lines) {
      if (line.includes('Session manager initialized')) {
        this.emit('session-manager-initialized')
        continue
      }

      if (line.includes('connect to host') && line.includes('Connection refused')) {
        this.emit('connection-to-server-dropped')
        continue
      }

      const match = line.match(MutagenDaemon.#MUTAGEN_STATE_TRANSITION_LOG_PATTERN)
      if (match == null) {
        continue
      }
      const [, sessionName, fromState, toState] = match

      this.emit(
        'session-state-changed',
        sessionName,
        fromState as MutagenSessionState,
        toState as MutagenSessionState
      )
    }
  }

  async #handleStdio(stream: Readable) {
    for await (const chunk of stream) {
      this.#onStdioChunk(chunk)
      this.#executable.log.write(chunk)
    }
  }

  async onRunning() {
    if (!this.isRunning) {
      await new Promise<void>((resolve) => this.once('session-manager-initialized', resolve))
    }
  }

  async start() {
    if (this.#process != null) {
      await this.onRunning()
      return
    }

    const daemonProcess = (this.#process = this.#executable.spawn(['daemon', 'run'], {
      stdio: ['ignore', 'pipe', 'pipe'],
    }))

    this.#handleStdio(daemonProcess.stdout)
    this.#handleStdio(daemonProcess.stderr)

    await new Promise<void>((resolve, reject) => {
      const daemon = this

      daemon.once('session-manager-initialized', onStarted)
      daemonProcess.once('exit', onExit)

      function onStarted() {
        daemonProcess.off('exit', onExit)
        resolve()
      }

      function onExit(status: number) {
        daemon.off('session-manager-initialized', onStarted)
        const error = new Error(`Failed to start daemon. Status code: ${status}`)
        daemon.emit('failed-to-start', error)
        daemon.#process = undefined
        reject(error)
      }
    })

    this.emit('is-running-changed', true)
    this.#process.on('exit', () => {
      this.#process = undefined
      this.emit('is-running-changed', false)
    })
  }

  async restart() {
    await this.stop()
    await this.start()
  }

  async stop() {
    await this.#kill(15) // SIGTERM
  }

  async kill() {
    await this.#kill(9) // SIGKILL
  }

  async #kill(signal: number) {
    if (this.#process == null) {
      return
    }
    const process = this.#process
    process.kill(signal)
    await new Promise((resolve) => process.once('exit', resolve))
  }

  async deleteDir() {
    // Length is suspiciously short, don't trigger delete (just in case)
    if (this.#executable.dataDirectory.length < 10) {
      throw new Error('not deleting mutagen dir, path is very short')
    }
    this.#logger.log('deleting mutagen dir', this.#executable.dataDirectory)
    await rm(this.#executable.dataDirectory, { recursive: true, force: true })
  }
}

import { Client, gql } from '@urql/core'
import path from 'path'
import { appendFile, mkdir, readFile, stat, rename, chmod } from 'fs/promises'
import { createWriteStream } from 'fs'
import { spawn } from 'child_process'
import { homedir } from 'os'
import { MessageChannel, Worker } from 'worker_threads'
import { Status } from './application'
import { Logger } from './Logger'

const keyMode = 0o600
const dotSSHMode = 0o755
const knownHostsMode = 0o644

const ensureMode = async (filepath: string, mode: number) =>
  await stat(filepath).then((stat) => {
    if (stat.mode !== mode) {
      return chmod(filepath, mode)
    }
  })

export class SSHKeys {
  readonly #logger: Logger
  readonly #client: Client
  readonly #status: Status
  readonly #syncHostURL: URL
  readonly #directory: string
  readonly #worker: Worker
  #sequence: Promise<any> = Promise.resolve()

  constructor(logger: Logger, client: Client, status: Status, syncHostURL: URL, directory: string) {
    this.#logger = logger.withPrefix('ssh-keys')
    this.#client = client
    this.#status = status
    this.#syncHostURL = syncHostURL
    this.#directory = directory
    this.#worker = new Worker(new URL('./sshWorker.js', import.meta.url))
  }

  ensure(): Promise<string> {
    const promise = this.#sequence.then(async () => {
      const { data, error } = await this.#client
        .query(
          gql`
            {
              user {
                id
              }
            }
          `
        )
        .toPromise()

      if (error != null) {
        throw error
      }

      const userID: string | undefined = data?.user?.id

      if (userID == null) {
        throw new Error('Authentication failed')
      }

      await this.#trustSyncHost()

      const keyPath = this.#path(userID)

      try {
        const { size } = await stat(keyPath)
        if (size === 0) {
          await this.#create(userID)
        } else {
          await ensureMode(keyPath, keyMode).catch((e) => {
            this.#logger.error(`Failed to set mode on ${keyPath}: ${e}`)
          })
        }
        return keyPath
      } catch (e) {
        if ((e as NodeJS.ErrnoException).code === 'ENOENT') {
          await this.#create(userID)
          return keyPath
        }
        throw e
      }
    })

    this.#sequence = promise

    return promise
  }

  #path(userID: string): string {
    return path.join(this.#directory, `private-key-ed25519-${userID}.pem`)
  }

  async #create(userID: string) {
    this.#status.willCreateSSHKeys()
    const path = this.#path(userID)
    const tmpPath = path + '.tmp'

    // Will throw if the directory doesn't exist.
    const stream = createWriteStream(tmpPath, {
      encoding: 'ascii',
      mode: keyMode,
    })

    const { port1: callback, port2: replyTo } = new MessageChannel()

    this.#worker.postMessage(replyTo, [replyTo])

    const [privateKey, publicKey] = await new Promise<[string, string]>((r) =>
      callback.once('message', r)
    )
    this.#status.willUploadSSHKeys()

    await this.#uploadPublicKey(publicKey)

    await new Promise<void>((r) => stream.end(privateKey, r))

    await rename(tmpPath, path)

    this.#status.createdAndUploadedSSHKeys()
  }

  async #uploadPublicKey(publicKey: string) {
    const { error } = await this.#client
      .mutation(
        gql`
          mutation ($publicKey: String!) {
            addPublicKey(publicKey: $publicKey) {
              id
            }
          }
        `,
        {
          publicKey,
        }
      )
      .toPromise()

    if (error != null) {
      throw error
    }
  }

  async #trustSyncHost() {
    this.#logger.log('Adding sync trust')

    const keyscan = spawn(
      'ssh-keyscan',
      ['-p', this.#syncHostURL.port || '22', this.#syncHostURL.hostname],
      {
        stdio: ['ignore', 'pipe', 'pipe'],
      }
    )
    const chunks: Buffer[] = []
    keyscan.stdout.on('data', (chunk) => {
      chunks.push(chunk)
      console.log('ssh-keyscan stdout', chunk.toString())
    })
    keyscan.stderr.on('data', (chunk) => {
      console.log('ssh-keyscan stderr', chunk.toString())
    })

    const statusCode = await new Promise<number>((resolve, reject) =>
      keyscan.once('error', reject).once('exit', resolve)
    )

    if (statusCode !== 0) {
      throw new Error(`Failed to add trust. Status code: ${statusCode}`)
    }

    const trustRows = Buffer.concat(chunks)
      .toString('utf-8')
      .split(/[\r\n]+/g)
      .filter((r) => !r.startsWith('#'))

    const sshDir = path.join(homedir(), '.ssh')
    const knownHostsFile = path.join(sshDir, 'known_hosts')

    // Create SSH directory if not exists
    await mkdir(sshDir, {
      mode: dotSSHMode,
      recursive: true,
    })

    await ensureMode(sshDir, dotSSHMode).catch((e) => {
      this.#logger.error(`Failed to set mode on ${sshDir}: ${e}`)
    })

    let existingTrustRows = ''
    try {
      existingTrustRows = await readFile(knownHostsFile, 'utf-8')
    } catch {}

    const newTrustRows: string[] = []
    for (const row of trustRows) {
      if (!existingTrustRows.includes(row.trim())) {
        newTrustRows.push(row)
      }
    }

    if (newTrustRows.length === 0) {
      return
    }

    const lineSeparator = process.platform === 'win32' ? '\r\n' : '\n'

    await appendFile(knownHostsFile, newTrustRows.map((r) => r + lineSeparator).join(''), {
      mode: knownHostsMode,
    })
    await ensureMode(knownHostsFile, knownHostsMode).catch((e) => {
      this.#logger.error(`Failed to set mode on ${knownHostsFile}: ${e}`)
    })
  }
}

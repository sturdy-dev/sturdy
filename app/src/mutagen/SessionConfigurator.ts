import { Client, gql } from '@urql/core'
import { SSHKeys } from '../SSHKeys'
import { writeFile, mkdir } from 'fs/promises'
import { join as joinPath } from 'path'
import { MutagenSession, SessionConfig } from './Session'
import { MutagenExecutable } from './Executable'
import { MutagenDaemon } from './Daemon'
import { Logger } from '../Logger'

interface MutagenConfig {
  sync: MutagenSyncConfig
}

interface MutagenSyncConfig {
  defaults: MutagenSyncEndpointConfig
}

interface MutagenSyncEndpointConfig {
  mode: string
  sshPrivateKeyPath: string
  ignore: MutagenSyncEndpointIgnoreConfig
}

interface MutagenSyncEndpointIgnoreConfig {
  paths: string[]
  vcs: boolean
}

export class MutagenSessionConfigurator {
  readonly #logger: Logger
  readonly #directory: string
  readonly #executable: MutagenExecutable
  readonly #daemon: MutagenDaemon
  readonly #sshKeys: SSHKeys
  readonly #apiURL: URL
  readonly #syncHostURL: URL
  readonly #client: Client

  constructor(
    logger: Logger,
    directory: string,
    executable: MutagenExecutable,
    daemon: MutagenDaemon,
    sshKeys: SSHKeys,
    apiURL: URL,
    syncHostURL: URL,
    client: Client
  ) {
    this.#logger = logger.withPrefix('session')
    this.#directory = directory
    this.#executable = executable
    this.#daemon = daemon
    this.#sshKeys = sshKeys
    this.#apiURL = apiURL
    this.#syncHostURL = syncHostURL
    this.#client = client
  }

  get executable() {
    return this.#executable
  }

  get daemon() {
    return this.#daemon
  }

  get apiURL() {
    return this.#apiURL
  }

  async configureAndStart(viewID: string, path: string): Promise<MutagenSession> {
    return MutagenSession.create(await this.configure(viewID, path))
  }

  async configure(viewID: string, path: string): Promise<SessionConfig> {
    const { data, error } = await this.#client
      .query(
        gql`
          query ($viewID: ID!) {
            user {
              id
            }
            view(id: $viewID) {
              ignoredPaths
              codebase {
                id
              }
            }
          }
        `,
        { viewID }
      )
      .toPromise()

    if (error != null) {
      throw error
    }

    const userID = data?.user?.id
    const codebaseID = data?.view?.codebase?.id

    if (userID == null || codebaseID == null) {
      throw new Error("View doesn't exist")
    }

    const privateKeyPath = await this.#sshKeys.ensure()

    const ignores = new Set([
      ...(data.view.ignoredPaths ?? []),
      'node_modules',
      '.DS_Store',
      '*.swp',
    ])

    const conf: MutagenConfig = {
      sync: {
        defaults: {
          mode: 'two-way-resolved',
          sshPrivateKeyPath: privateKeyPath,
          ignore: {
            paths: Array.from(ignores),
            vcs: true,
          },
        },
      },
    }

    const sessionConfigPath = joinPath(this.#directory, `${viewID}.yaml`)

    await writeFile(sessionConfigPath, JSON.stringify(conf))

    await mkdir(path, { recursive: true })

    return {
      logger: this.#logger,
      configPath: sessionConfigPath,
      mountPath: path,
      viewID,
      codebaseID,
      userID,
      apiURL: this.#apiURL,
      syncHostURL: this.#syncHostURL,
      executable: this.#executable,
      daemon: this.#daemon,
    }
  }
}

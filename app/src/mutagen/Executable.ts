import {
  ChildProcess,
  ChildProcessByStdio,
  spawn,
  SpawnOptions,
  SpawnOptionsWithStdioTuple,
  StdioNull,
  StdioOptions,
  StdioPipe,
} from 'child_process'
import { Readable, Writable } from 'node:stream'
import { Logger } from '../Logger'

export class MutagenExecutable {
  readonly #executablePath: string
  readonly #runningProcesses = new Set<ChildProcess>()
  readonly #logger: Logger

  constructor({ executablePath, logger }: { executablePath: string; logger: Logger }) {
    this.#executablePath = executablePath
    this.#logger = logger.withPrefix('mutagen')
  }

  #decorateSpawnOptions<O extends SpawnOptions>(
    options: O,
    log: Writable,
    dataDirectory: string
  ): O {
    let stdio: StdioOptions
    switch (typeof options.stdio) {
      case 'undefined':
        stdio = ['ignore', log, log]
        break
      case 'string':
        stdio = options.stdio === 'ignore' ? ['ignore', log, log] : options.stdio
        break
      default:
        stdio = options.stdio.map((io, idx) => (idx > 0 && io === 'ignore' ? log : io))
        break
    }
    return {
      ...options,
      env: {
        MUTAGEN_DATA_DIRECTORY: dataDirectory,
        MUTAGEN_DISABLE_AUTOSTART: '1',
        ...(options.env ?? process.env),
      },
      stdio,
    }
  }

  abort() {
    if (this.#runningProcesses.size > 0) {
      this.#logger.log('killing', this.#runningProcesses.size, 'in-flight sturdy-sync processes')
      for (const process of this.#runningProcesses) {
        process.kill()
      }
    }
  }

  execute(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioPipe, StdioPipe, StdioPipe>,
    log: Writable,
    dataDirectory: string
  ): [proc: ChildProcessByStdio<Writable, Readable, Readable>, onExit: Promise<void>]
  execute(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioPipe, StdioPipe, StdioNull>,
    log: Writable,
    dataDirectory: string
  ): [proc: ChildProcessByStdio<Writable, Readable, null>, onExit: Promise<void>]
  execute(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioPipe, StdioNull, StdioPipe>,
    log: Writable,
    dataDirectory: string
  ): [proc: ChildProcessByStdio<Writable, null, Readable>, onExit: Promise<void>]
  execute(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioNull, StdioPipe, StdioPipe>,
    log: Writable,
    dataDirectory: string
  ): [proc: ChildProcessByStdio<null, Readable, Readable>, onExit: Promise<void>]
  execute(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioPipe, StdioNull, StdioNull>,
    log: Writable,
    dataDirectory: string
  ): [proc: ChildProcessByStdio<Writable, null, null>, onExit: Promise<void>]
  execute(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioNull, StdioPipe, StdioNull>,
    log: Writable,
    dataDirectory: string
  ): [proc: ChildProcessByStdio<null, Readable, null>, onExit: Promise<void>]
  execute(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioNull, StdioNull, StdioPipe>,
    log: Writable,
    dataDirectory: string
  ): [proc: ChildProcessByStdio<null, null, Readable>, onExit: Promise<void>]
  execute(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioNull, StdioNull, StdioNull>,
    log: Writable,
    dataDirectory: string
  ): [proc: ChildProcessByStdio<null, null, null>, onExit: Promise<void>]
  execute(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<any, any, any>,
    log: Writable,
    dataDirectory: string
  ): [proc: ChildProcess, onExit: Promise<void>] {
    const proc = spawn(
      this.#executablePath,
      args,
      this.#decorateSpawnOptions(options, log, dataDirectory)
    )

    this.#register(proc)

    const onExit = new Promise<void>((resolve, reject) => {
      proc.once('exit', (status) => {
        if (status !== 0) {
          reject(
            new Error(
              `Command ${this.#executablePath} ${args.join(' ')} failed. Status code: ${status}`
            )
          )
        } else {
          resolve()
        }
      })
    })

    return [proc, onExit]
  }

  spawn(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioPipe, StdioPipe, StdioPipe>,
    log: Writable,
    dataDirectory: string
  ): ChildProcessByStdio<Writable, Readable, Readable>
  spawn(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioPipe, StdioPipe, StdioNull>,
    log: Writable,
    dataDirectory: string
  ): ChildProcessByStdio<Writable, Readable, null>
  spawn(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioPipe, StdioNull, StdioPipe>,
    log: Writable,
    dataDirectory: string
  ): ChildProcessByStdio<Writable, null, Readable>
  spawn(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioNull, StdioPipe, StdioPipe>,
    log: Writable,
    dataDirectory: string
  ): ChildProcessByStdio<null, Readable, Readable>
  spawn(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioPipe, StdioNull, StdioNull>,
    log: Writable,
    dataDirectory: string
  ): ChildProcessByStdio<Writable, null, null>
  spawn(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioNull, StdioPipe, StdioNull>,
    log: Writable,
    dataDirectory: string
  ): ChildProcessByStdio<null, Readable, null>
  spawn(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioNull, StdioNull, StdioPipe>,
    log: Writable,
    dataDirectory: string
  ): ChildProcessByStdio<null, null, Readable>
  spawn(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioNull, StdioNull, StdioNull>,
    log: Writable,
    dataDirectory: string
  ): ChildProcessByStdio<null, null, null>
  spawn(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<any, any, any>,
    log: Writable,
    dataDirectory: string
  ): ChildProcess {
    const spawnOpts = this.#decorateSpawnOptions(options, log, dataDirectory)
    this.#logger.log('spawn', {
      path: this.#executablePath,
      args,
      cwd: spawnOpts.cwd,
      dataDirectory,
    })
    const process = spawn(this.#executablePath, args, spawnOpts)
    this.#register(process)
    return process
  }

  #register(process: ChildProcess) {
    this.#runningProcesses.add(process)
    process.once('exit', () => {
      this.#runningProcesses.delete(process)
    })
  }
}

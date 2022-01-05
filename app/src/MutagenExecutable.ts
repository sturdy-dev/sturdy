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

export class MutagenExecutable {
  readonly #executablePath: string
  readonly #dataDirectory: string
  readonly #runningProcesses = new Set<ChildProcess>()
  readonly #log: Writable

  constructor({
    executablePath,
    dataDirectory,
    log,
  }: {
    executablePath: string
    dataDirectory: string
    log: Writable
  }) {
    this.#executablePath = executablePath
    this.#dataDirectory = dataDirectory
    this.#log = log
  }

  #decorateSpawnOptions<O extends SpawnOptions>(options: O): O {
    let stdio: StdioOptions
    switch (typeof options.stdio) {
      case 'undefined':
        stdio = ['ignore', this.#log, this.#log]
        break
      case 'string':
        stdio = options.stdio === 'ignore' ? ['ignore', this.#log, this.#log] : options.stdio
        break
      default:
        stdio = options.stdio.map((io, idx) => (idx > 0 && io === 'ignore' ? this.#log : io))
        break
    }
    return {
      ...options,
      env: {
        MUTAGEN_DATA_DIRECTORY: this.#dataDirectory,
        MUTAGEN_DISABLE_AUTOSTART: '1',
        ...(options.env ?? process.env),
      },
      stdio,
    }
  }

  abort() {
    if (this.#runningProcesses.size > 0) {
      console.log('killing', this.#runningProcesses.size, 'in-flight sturdy-sync processes')
      for (const process of this.#runningProcesses) {
        process.kill()
      }
    }
  }

  execute(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioPipe, StdioPipe, StdioPipe>
  ): [proc: ChildProcessByStdio<Writable, Readable, Readable>, onExit: Promise<void>]
  execute(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioPipe, StdioPipe, StdioNull>
  ): [proc: ChildProcessByStdio<Writable, Readable, null>, onExit: Promise<void>]
  execute(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioPipe, StdioNull, StdioPipe>
  ): [proc: ChildProcessByStdio<Writable, null, Readable>, onExit: Promise<void>]
  execute(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioNull, StdioPipe, StdioPipe>
  ): [proc: ChildProcessByStdio<null, Readable, Readable>, onExit: Promise<void>]
  execute(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioPipe, StdioNull, StdioNull>
  ): [proc: ChildProcessByStdio<Writable, null, null>, onExit: Promise<void>]
  execute(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioNull, StdioPipe, StdioNull>
  ): [proc: ChildProcessByStdio<null, Readable, null>, onExit: Promise<void>]
  execute(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioNull, StdioNull, StdioPipe>
  ): [proc: ChildProcessByStdio<null, null, Readable>, onExit: Promise<void>]
  execute(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioNull, StdioNull, StdioNull>
  ): [proc: ChildProcessByStdio<null, null, null>, onExit: Promise<void>]
  execute(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<any, any, any>
  ): [proc: ChildProcess, onExit: Promise<void>] {
    const proc = spawn(this.#executablePath, args, this.#decorateSpawnOptions(options))

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
    options: SpawnOptionsWithStdioTuple<StdioPipe, StdioPipe, StdioPipe>
  ): ChildProcessByStdio<Writable, Readable, Readable>
  spawn(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioPipe, StdioPipe, StdioNull>
  ): ChildProcessByStdio<Writable, Readable, null>
  spawn(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioPipe, StdioNull, StdioPipe>
  ): ChildProcessByStdio<Writable, null, Readable>
  spawn(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioNull, StdioPipe, StdioPipe>
  ): ChildProcessByStdio<null, Readable, Readable>
  spawn(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioPipe, StdioNull, StdioNull>
  ): ChildProcessByStdio<Writable, null, null>
  spawn(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioNull, StdioPipe, StdioNull>
  ): ChildProcessByStdio<null, Readable, null>
  spawn(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioNull, StdioNull, StdioPipe>
  ): ChildProcessByStdio<null, null, Readable>
  spawn(
    args: readonly string[],
    options: SpawnOptionsWithStdioTuple<StdioNull, StdioNull, StdioNull>
  ): ChildProcessByStdio<null, null, null>
  spawn(args: readonly string[], options: SpawnOptionsWithStdioTuple<any, any, any>): ChildProcess {
    const spawnOpts = this.#decorateSpawnOptions(options)
    console.log('spawn', {
      path: this.#executablePath,
      args,
      cwd: spawnOpts.cwd,
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

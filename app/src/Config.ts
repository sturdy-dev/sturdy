export interface ViewConfig {
  id: string
  path: string
}

export interface Config {
  views: ViewConfig[]
}

export namespace Config {
  export function defaultConfig(): Config {
    return {
      views: [],
    }
  }
}

export interface Crumb {
  name: string
  fullPath: string
  isCurrent: boolean
}

export class Breadcrumbs {
  static paths(fullPath: string): Crumb[] {
    const crumbs = Array<Crumb>()

    if (fullPath === '/' || fullPath === '') {
      return crumbs
    }

    let upTo = ''
    const parts = fullPath.split('/')

    for (const [idx, part] of parts.entries()) {
      upTo += part
      crumbs.push({
        name: part,
        fullPath: upTo,
        isCurrent: idx >= parts.length - 1,
      })
      upTo += '/'
    }

    return crumbs
  }
}

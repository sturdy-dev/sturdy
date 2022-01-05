import { app } from 'electron'
import path from 'path'

export function resourcePath(...segments: string[]) {
  if (app.isPackaged) {
    return path.join(process.resourcesPath, 'app', 'assets', ...segments)
  }
  return path.join(process.cwd(), 'assets', ...segments)
}

export function dataPath(...segments: string[]) {
  return path.join(app.getPath('userData'), ...segments)
}

export function resourceBinary(...segments: string[]) {
  if (process.platform === 'win32') {
    return resourcePath('bin', ...segments) + '.exe'
  }
  return resourcePath('bin', ...segments)
}

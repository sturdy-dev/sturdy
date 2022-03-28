const https = require('https')
const fs = require('fs')
const path = require('path')
const tar = require('tar')
const unzipper = require('unzipper')

const STURDY_SYNC_VERSION = 'v0.9.0'

const log = (message) => console.log(`  * ${message}`)
const error = (message) => console.error(`  * ${message}`)

const downloadBinaries = async (platform, arch, dest) => {
  fs.mkdirSync(dest, { recursive: true })

  const format = platform === 'windows' ? 'zip' : 'tar.gz'
  const archiveName = `sturdy-${STURDY_SYNC_VERSION}-${platform}-${arch}.${format}`
  const downloadUrl = `https://getsturdy.com/client/${archiveName}`

  log(`downloading ${downloadUrl}`)

  await new Promise((resolve, reject) => {
    https.get(downloadUrl, (res) => {
      const unarchive =
        platform === 'windows' ? unzipper.Extract({ path: dest }) : tar.x({ cwd: dest })
      unarchive.on('error', reject)

      res.pipe(unarchive)
      res.on('end', resolve)
      res.on('error', reject)
    })
  }).catch((err) => {
    error(`failed to download`)
    throw err
  })

  fs.readdirSync(dest)
    .map((file) => `downloaded: ${path.join(dest, file)}`)
    .forEach(log)
}

const translatePlatform = (platform) => {
  switch (platform) {
    case 'windows':
      return 'windows'
    case 'mac':
      return 'darwin'
    case 'linux':
      return 'linux'
    default:
      throw new Error(`Unsupported platform: ${platform}`)
  }
}

const translateArch = (arch) => {
  switch (arch) {
    case 1:
      return 'amd64'
    case 3:
      return 'arm64'
    default:
      throw new Error(`Unsupported arch: ${arch}`)
  }
}

/** @type {import('electron-builder').Configuration['beforePack']} */
module.exports = async function (params) {
  const platform = translatePlatform(params.packager.platform.name)
  const arch = translateArch(params.arch)
  const downloadTo = path.join(params.appOutDir, 'resources/app/assets/bin')
  await downloadBinaries(platform, arch, downloadTo)
}

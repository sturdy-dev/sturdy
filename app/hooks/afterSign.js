// See: https://medium.com/@TwitterArchiveEraser/notarize-electron-apps-7a5f988406db

const fs = require('fs')
const path = require('path')
const electron_notarize = require('electron-notarize')

const log = (message) => {
  console.log(`  * ${message}`)
}

const error = (message) => {
  console.error(`  * ${message}`)
}

/** @type {import('electron-builder').Configuration['afterSign']} */
module.exports = async function (params) {
  // Only sign mac builds
  if (params.packager.platform.name !== 'mac') {
    return
  }

  if (process.platform !== 'darwin') {
    throw new Error('Notarization is only supported on MacOS')
  }

  if (!process.env.APPLE_ID) {
    log('No APPLE_ID environment variable found. Skipping notarization.')
    return
  }

  const appId = params.packager.appInfo.id
  const appPath = path.join(params.appOutDir, `${params.packager.appInfo.productFilename}.app`)
  if (!fs.existsSync(appPath)) {
    throw new Error(`Cannot find application at: ${appPath}`)
  }

  log(`notarizing ${appPath}`)

  await electron_notarize
    .notarize({
      appBundleId: appId,
      appPath: appPath,
      appleId: process.env.APPLE_ID,
      appleIdPassword: process.env.APPLE_ID_PASSWORD,
    })
    .catch((err) => {
      error(`notarization failed: ${err.message}`)
      throw err
    })
    .finally(() => {
      log('notarization complete')
    })
}

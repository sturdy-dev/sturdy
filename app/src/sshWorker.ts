import { MessagePort, parentPort } from 'worker_threads'

const { generateKeyPairSync } = require('crypto')
const sshpk = require('sshpk')

parentPort?.on('message', (replyTo: MessagePort) => {
  // Generate a new keypair
  const { privateKey } = generateKeyPairSync('ed25519')
  let privateStr = privateKey.export({ format: 'pem', type: 'pkcs8' })

  // Format it in OpenSSH format. The public key will be derived from the private key.
  const sshPublicKey = sshpk.parseKey(privateStr, 'pem').toString('ssh')
  const sshPrivateKey = sshpk.parsePrivateKey(privateStr, 'pem').toString('ssh')

  replyTo.postMessage([sshPrivateKey, sshPublicKey])
})

import typescript from '@rollup/plugin-typescript'
import { nodeResolve } from '@rollup/plugin-node-resolve'
import packageJson from './package.json'

export default {
  input: {
    main: './src/main.ts',
    preload: './src/preload.ts',
    sshWorker: './src/sshWorker.ts',
  },
  external: [...Object.keys(packageJson.dependencies), 'electron', 'fs/promises', 'node-fetch'],
  output: {
    dir: 'dist',
    format: 'cjs',
  },
  plugins: [typescript(), nodeResolve(), retainImportExpressionPlugin()],
}

function retainImportExpressionPlugin() {
  return {
    name: 'retain-import-expression',
    resolveDynamicImport(specifier) {
      if (specifier === 'node-fetch') return false
      return null
    },
    renderDynamicImport({ targetModuleId }) {
      if (targetModuleId === 'node-fetch') {
        return {
          left: 'import(',
          right: ')',
        }
      }
    },
  }
}

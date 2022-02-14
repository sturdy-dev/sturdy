import typescript from '@rollup/plugin-typescript'
import resolve from '@rollup/plugin-node-resolve'
import packageJson from './package.json'
import preferencesConfig from './preferences.rollup.config.js'

export default [
  {
    input: {
      main: './src/main.ts',
      preload: './src/preload.ts',
      sshWorker: './src/sshWorker.ts',
      'preferences/preload': './src/preferences/preload.ts',
    },
    external: [...Object.keys(packageJson.dependencies), 'electron', 'fs/promises', 'node-fetch'],
    output: {
      dir: 'dist',
      format: 'cjs',
    },
    plugins: [
      typescript({
        tsconfig: './tsconfig.json',
      }),
      resolve(),
      retainImportExpressionPlugin(),
    ],
  },
  preferencesConfig,
]

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

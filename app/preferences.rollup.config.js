import vue from 'rollup-plugin-vue'
import typescript from '@rollup/plugin-typescript'
import resolve from '@rollup/plugin-node-resolve'
import postcss from 'rollup-plugin-postcss'
import injectProcessEnv from 'rollup-plugin-inject-process-env'
import html from '@rollup/plugin-html'

export default {
  input: './src/preferences/app/main.ts',
  output: {
    dir: 'dist/preferences',
    format: 'es',
  },
  plugins: [
    typescript({
      tsconfig: false,
      include: ['src/preferences/app/**/*.ts', 'src/preferences/app/**/*.vue'],
      exclude: ['node_modules'],
      moduleResolution: 'node',
      target: 'esnext',
      module: 'esnext',
    }),
    vue(),
    resolve(),
    postcss({
      config: {
        path: 'postcss.config.js',
      },
    }),
    injectProcessEnv({
      NODE_ENV: 'production',
    }),
    html({
      title: 'Preferences',
    }),
  ],
}

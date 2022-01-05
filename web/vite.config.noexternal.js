/* eslint-disable @typescript-eslint/no-var-requires */

const config = require('./vite.config.js')
/**
 * @type {import('vite').UserConfig}
 */
module.exports = Object.assign(config, {
  ssr: {
    noExternal: /heroicons/,
    target: 'node',
  },
  resolve: {
    // necessary because vue.ssrUtils is only exported on cjs modules
    alias: [
      {
        find: '@vue/runtime-dom',
        replacement: '@vue/runtime-dom/dist/runtime-dom.cjs.js',
      },
      {
        find: '@vue/runtime-core',
        replacement: '@vue/runtime-core/dist/runtime-core.cjs.js',
      },
    ],
  },
  build: {
    minify: false,
    // rollupOptions is not set
  },
})

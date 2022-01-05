/* eslint-disable @typescript-eslint/no-var-requires */

const vuePlugin = require('@vitejs/plugin-vue')
const vueJsx = require('@vitejs/plugin-vue-jsx')
const { resolve } = require('path')
import { imagetools } from 'vite-imagetools'
import eslintPlugin from 'vite-plugin-eslint'

const ssrTransformCustomDir = () => {
  return {
    props: [],
    needRuntime: true,
  }
}

/**
 * @type {import('vite').UserConfig}
 */
module.exports = {
  plugins: [
    imagetools(),
    eslintPlugin({ cache: false }),
    vuePlugin({
      template: {
        ssr: true,
        compilerOptions: {
          directiveTransforms: {
            'click-outside': ssrTransformCustomDir,
            'mousedown-outside': ssrTransformCustomDir,
          },
        },
      },
    }),
    vueJsx(),
    {
      name: 'virtual',
      resolveId(id) {
        if (id === '@foo') {
          return id
        }
      },
      load(id) {
        if (id === '@foo') {
          return `export default { msg: 'hi' }`
        }
      },
    },
  ],
  build: {
    minify: false,
    rollupOptions: {
      input: {
        main: resolve(__dirname, 'index.html'),
        nested: resolve(__dirname, 'client-side-render.html'),
      },
    },
  },
  server: {
    port: 8080,
  },
}

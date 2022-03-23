import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import { imagetools } from 'vite-imagetools'
import eslintPlugin from 'vite-plugin-eslint'

const ssrTransformCustomDir = () => {
  return {
    props: [],
    needRuntime: true,
  }
}

const isProd = (process.env.VITE_ENV = 'production')

/**
 * @type {import('vite').UserConfig}
 */
module.exports = {
  plugins: [
    imagetools(),
    eslintPlugin({ cache: false }),
    vue({
      isProduction: isProd,
      template: {
        isProd,
        ssr: true,
        compilerOptions: {
          directiveTransforms: {
            'click-outside': ssrTransformCustomDir,
            'mousedown-outside': ssrTransformCustomDir,
          },
        },
      },
    }),
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

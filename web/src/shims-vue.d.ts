declare module '*.vue' {
  import { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

declare namespace JSX {
  import * as dom from '@vue/runtime-dom'

  type IntrinsicElements = dom.IntrinsicElementAttributes
}

import mitt from 'mitt'
import { Feature } from './features'

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties {
    emitter: mitt
    features: Set<Feature>
  }
}

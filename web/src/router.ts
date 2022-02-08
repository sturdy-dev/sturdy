import {
  createMemoryHistory,
  createRouter,
  createWebHashHistory,
  createWebHistory,
  RouterOptions,
} from 'vue-router'
import { RoutesDocs } from './router-docs'
import { RoutesApps, RoutesSelfHosted } from './router-apps'

const buildDocs =
  import.meta.env.VITE_WEB_BUILD_DOCS && import.meta.env.VITE_WEB_BUILD_DOCS === 'true'

const routes = (buildDocs ? RoutesDocs : RoutesSelfHosted).concat(RoutesApps)

const routerOpts: Partial<RouterOptions> = {
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else if (to.hash) {
      return { el: to.hash, scrollBehavior: 'smooth', top: 150 }
    } else {
      return { top: 0 }
    }
  },
}

if (import.meta.env.SSR) {
  routerOpts.history = createMemoryHistory()
} else if (import.meta.env.VITE_USE_HASH_HISTORY) {
  routerOpts.history = createWebHashHistory()
} else {
  routerOpts.history = createWebHistory()
}

export default createRouter(routerOpts as RouterOptions)

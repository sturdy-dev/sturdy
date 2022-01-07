import { createApp as vueCreateApp, createSSRApp } from 'vue'
import Sturdy from './Sturdy.vue'
import router from './router'
import mitt from 'mitt'
import './index.css'
import urql, { dedupExchange, fetchExchange, ssrExchange, subscriptionExchange } from '@urql/vue'
import { retryExchange } from '@urql/exchange-retry'
import { SubscriptionClient } from 'subscriptions-transport-ws'
import { devtoolsExchange } from '@urql/devtools'
import { createHead } from '@vueuse/head'
import { RetryExchangeOptions } from '@urql/exchange-retry/dist/types/retryExchange'
import { cacheExchange } from '@urql/exchange-graphcache'
import schema from '../schema.json'
import { IntrospectionData } from '@urql/exchange-graphcache/dist/types/ast'
import { subscriptionUpdateResolvers } from './subscriptions/subscriptionUpdateResolvers'
import {
  mutationUpdateResolvers,
  optimisticMutationResolvers,
} from './mutations/mutationUpdateResolvers'
import { keyResolvers } from './keys/keyResolvers'

export function createApp(ssrApp: boolean) {
  // Global message bus
  const emitter = mitt()

  let app

  if (ssrApp) {
    app = createSSRApp(Sturdy)
  } else {
    app = vueCreateApp(Sturdy)
  }

  app.use(router)

  const head = createHead()
  app.use(head)

  app.config.globalProperties.emitter = emitter

  const options: RetryExchangeOptions = {
    initialDelayMs: 1000,
    maxDelayMs: 15000,
    randomDelay: true,
    maxNumberAttempts: 2,
    retryIf: (err) =>
      Boolean(err && err.networkError && err.networkError.message !== 'Unauthorized'),
  }

  const exchanges = [
    dedupExchange,
    cacheExchange({
      schema: schema as IntrospectionData,
      updates: {
        Subscription: subscriptionUpdateResolvers,
        Mutation: mutationUpdateResolvers,
      },
      optimistic: optimisticMutationResolvers,
      keys: keyResolvers,
    }),
    ssrExchange({
      isClient: !import.meta.env.SSR,
      // initialState: !import.meta.env.SSR ? window.__URQL_DATA__ : undefined,
    }),
    retryExchange(options), // Use the retryExchange factory to add a new exchange
  ]

  // Client-side only
  if (!import.meta.env.SSR) {
    const wsHost = (import.meta.env.VITE_API_HOST as string)
      .replace('http://', 'ws://')
      .replace('https://', 'wss://')

    const subscriptionClient = new SubscriptionClient(wsHost + 'graphql/ws', {
      reconnect: true,
    })

    // Add client side only exchanges
    exchanges.push(
      fetchExchange,
      subscriptionExchange({
        forwardSubscription: (operation) => subscriptionClient.request(operation),
      })
    )
  }

  if (import.meta.env.VITE_ENABLE_URQL_DEVTOOLS) {
    exchanges.unshift(devtoolsExchange)
  }

  app.use(urql, {
    url: import.meta.env.VITE_API_HOST + 'graphql',
    fetchOptions: {
      credentials: 'include',
    },
    exchanges: exchanges,
  })

  return { app, router, head }
}

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
import { cacheExchange } from '@urql/exchange-graphcache'
import { getIntrospectedSchema, minifyIntrospectionQuery } from '@urql/introspection'
import { subscriptionUpdateResolvers } from './subscriptions/subscriptionUpdateResolvers'
import {
  mutationUpdateResolvers,
  optimisticMutationResolvers,
} from './mutations/mutationUpdateResolvers'
import { keyResolvers } from './keys/keyResolvers'
import schema from '../schema.json'

export function createApp(ssrApp: boolean) {
  // Global message bus
  const emitter = mitt()

  const app = ssrApp ? createSSRApp(Sturdy) : vueCreateApp(Sturdy)

  app.use(router)
  const head = createHead()
  app.use(head)

  app.config.globalProperties.emitter = emitter

  const exchanges = [
    dedupExchange,
    cacheExchange({
      schema: minifyIntrospectionQuery(getIntrospectedSchema(JSON.stringify(schema))),
      updates: {
        Subscription: subscriptionUpdateResolvers,
        Mutation: mutationUpdateResolvers,
      },
      optimistic: optimisticMutationResolvers,
      keys: keyResolvers,
    }),
    ssrExchange({
      isClient: !import.meta.env.SSR,
    }),
    retryExchange({
      initialDelayMs: 1000,
      maxDelayMs: 15000,
      randomDelay: true,
      maxNumberAttempts: 2,
      retryIf: (err) =>
        Boolean(err && err.networkError && err.networkError.message !== 'Unauthorized'),
    }), // Use the retryExchange factory to add a new exchange
  ]

  const host = import.meta.env.VITE_API_HOST
    ? (import.meta.env.VITE_API_HOST as string)
    : `${location.origin}`

  const apiPrefix = import.meta.env.VITE_API_PATH ?? ''

  // Client-side only
  if (!import.meta.env.SSR) {
    const wsHost = host.replace('http://', 'ws://').replace('https://', 'wss://')
    const graphqlWsUrl = `${wsHost}${apiPrefix}/graphql/ws`

    const subscriptionClient = new SubscriptionClient(graphqlWsUrl, {
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

  const graphqlUrl = `${host}${apiPrefix}/graphql`

  app.use(urql, {
    url: graphqlUrl,
    fetchOptions: {
      credentials: 'include',
    },
    exchanges: exchanges,
  })
  return { app, router, head }
}

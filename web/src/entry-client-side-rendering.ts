import { createApp } from './main'
import { getIntrospectionQuery, GraphQLSchema } from 'graphql'

// fetch the schema from the server
const fetchSchema = async (): Promise<GraphQLSchema> => {
  const resp = await fetch(`${import.meta.env.VITE_API_HOST}graphql/introspection`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      variables: {},
      query: getIntrospectionQuery({ descriptions: false }),
    }),
  })
  const { data } = await resp.json()
  return data
}

fetchSchema().then((schema) => {
  const { app } = createApp(false, schema)
  app.mount('#app')
})

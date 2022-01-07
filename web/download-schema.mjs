import { getIntrospectionQuery } from 'graphql'
import fetch from 'node-fetch' // or your preferred request in Node.js
import * as fs from 'fs'

import { getIntrospectedSchema, minifyIntrospectionQuery } from '@urql/introspection'

fetch('http://127.0.0.1:3000/graphql/introspection', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    variables: {},
    query: getIntrospectionQuery({ descriptions: false }),
  }),
})
  .then((result) => result.json())
  .then(({ data }) => {
    const minified = minifyIntrospectionQuery(getIntrospectedSchema(data))
    fs.writeFileSync('./schema.json', JSON.stringify(minified))
  })

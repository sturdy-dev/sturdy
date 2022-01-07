import { App } from 'vue'
import {
  IntrospectionField,
  IntrospectionObjectType,
  IntrospectionQuery,
  IntrospectionType,
} from 'graphql'

export type Feature = 'github'

export type ConfigOptions = {
  schema: IntrospectionQuery
}

const isGitHubEnabled = (query: IntrospectionQuery): boolean => {
  const schema = query.__schema
  if (!schema.types) return false
  const queries = schema.types.find(
    (type: IntrospectionType) => type.name === 'Query'
  ) as IntrospectionObjectType
  return queries.fields.some((field: IntrospectionField) => field.name === 'gitHubApp')
}

const buildFeatures = (introspection: IntrospectionQuery): Set<Feature> => {
  const set = new Set<Feature>()
  if (isGitHubEnabled(introspection)) set.add('github')
  return set
}

const install = (app: App, options: ConfigOptions) => {
  app.provide('features', buildFeatures(options.schema))
}

export default {
  install,
}

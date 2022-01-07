import { PartialIntrospectionSchema } from '@urql/exchange-graphcache/dist/types/ast'
import { App } from 'vue'
import {
  IntrospectionField,
  IntrospectionObjectType,
  IntrospectionSchema,
  IntrospectionType,
} from 'graphql'

export type Feature = 'github'

export type ConfigOptions = {
  schema: IntrospectionSchema
}

const isGitHubEnabled = (schema: IntrospectionSchema | PartialIntrospectionSchema): boolean => {
  if (!schema.types) return false
  const queries = schema.types.find(
    (type: IntrospectionType) => type.name === 'Query'
  ) as IntrospectionObjectType
  return queries.fields.some((field: IntrospectionField) => field.name === 'gitHubApp')
}

const buildFeatures = (introspection: IntrospectionSchema): Set<Feature> => {
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

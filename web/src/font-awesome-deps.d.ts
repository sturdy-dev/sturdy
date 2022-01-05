declare module '@fortawesome/fontawesome-svg-core' {
  import { IconDefinition } from '@fortawesome/fontawesome-common-types'

  export interface FontAwesomeLibrary {
    add(...definitions: IconDefinition[]): void
  }

  export const library: FontAwesomeLibrary
}

declare module '@fortawesome/free-solid-svg-icons' {
  import { IconDefinition } from '@fortawesome/fontawesome-common-types'

  export const fas: IconDefinition
}

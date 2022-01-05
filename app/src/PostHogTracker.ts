import { createClient, gql } from '@urql/core'
import PostHog from 'posthog-node'

const { app } = require('electron')

export class PostHogTracker {
  readonly #apiURL: URL
  #userID: string | undefined
  readonly #client: PostHog

  constructor(apiURL: URL, postHogToken: string) {
    this.#apiURL = apiURL
    this.#client = new PostHog(postHogToken)
  }

  unsetUser() {
    this.#userID = undefined
  }

  async updateUser(jwt: string): Promise<string> {
    const client = createClient({
      url: this.#apiURL.href,
      fetch: (await import('node-fetch')).default as any,
      fetchOptions: {
        credentials: 'include',
        headers: {
          Authorization: `bearer ${jwt}`,
        },
      },
    })

    const { data, error } = await client
      .query<{ user: { id: string; name: string } }>(
        gql`
          {
            user {
              id
              name
            }
          }
        `
      )
      .toPromise()

    if (error != null || data?.user == null) {
      return Promise.reject('could not get user: ' + error?.toString())
    }

    this.#userID = data.user.id

    return data.user.id
  }

  flush() {
    // On program exit, call shutdown to stop pending pollers and flush any remaining events
    this.#client.shutdown()
  }

  trackStartedApp() {
    if (!this.#userID) {
      return
    }

    this.#client.capture({
      distinctId: this.#userID,
      event: 'started app',
      properties: {
        // Send as event metadata
        'app-version': app.getVersion(),
        'app-platform': process.platform,

        // Also set last known app version and platform as user properties
        $set: {
          'app-version': app.getVersion(),
          'app-platform': process.platform,
        },
      },
    })
  }

  trackCreateNewView(viewID: string) {
    if (!this.#userID) {
      return
    }

    this.#client.capture({
      distinctId: this.#userID,
      event: 'create new view in app',
      properties: {
        viewID: viewID,
      },
    })
  }
}

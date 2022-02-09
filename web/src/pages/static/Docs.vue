<template>
  <StaticPage
    category="knowledgebase"
    title="Documentation"
    metadescription="Get to know Sturdy, and search for help"
    image="https://images.unsplash.com/photo-1531403009284-440f080d1e12?ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&ixlib=rb-1.2.1&auto=format&fit=crop&w=1050&q=80"
  >
    <div class="text-base max-w-prose mx-auto lg:max-w-none">
      <p class="text-lg text-gray-500">Learn how to use Sturdy and gain productivity superpowers</p>
    </div>
    <div
      class="mt-5 prose prose-yellow text-gray-500 mx-auto lg:max-w-none lg:row-start-1 lg:col-start-1"
    >
      <template v-for="[title, links] in groups" :key="title">
        <h2>{{ title }}</h2>

        <ul>
          <li v-for="link in links" :key="link.name">
            <router-link :to="{ name: link.name }" class="!no-underline">
              {{ link.meta.documentation.title }}
            </router-link>
          </li>
        </ul>
      </template>

      <h2>We're here to help!</h2>
      <p>
        Everyone at Sturdy is a developer, and we love to hear from other developers! We're here to
        help you with any problem that you might have <em>(sorry!)</em>, and listen to all of your
        cool ideas and feedback!
      </p>
    </div>

    <div class="mt-2">
      <ul class="space-y-2">
        <li>
          <LinkButton
            href="https://discord.com/invite/5HnSdzMqtA"
            target="_blank"
            class="!no-underline !text-black !m-0"
          >
            <img src="./docs/assets/discord.svg" class="h-8 w-8 mr-2" />
            <span>Get help on our Discord</span>
          </LinkButton>
        </li>
        <li>
          <LinkButton href="mailto:support@getsturdy.com" class="!no-underline !text-black !m-0">
            <AtSymbolIcon class="h-8 w-8 mr-2" />
            <span>support@getsturdy.com</span>
          </LinkButton>
        </li>
      </ul>
    </div>
  </StaticPage>
</template>

<script>
import StaticPage from '../../layouts/StaticPage.vue'
import { useRouter } from 'vue-router'
import LinkButton from '../../components/shared/LinkButton.vue'
import { AtSymbolIcon } from '@heroicons/vue/solid'

export default {
  components: { LinkButton, StaticPage, AtSymbolIcon },
  setup() {
    let routes = useRouter()
      .getRoutes()
      .filter((r) => r.meta.documentation)

    const groups = routes.reduce(
      (entryMap, e) =>
        entryMap.set(e.meta.documentation.group, [
          ...(entryMap.get(e.meta.documentation.group) || []),
          e,
        ]),
      new Map()
    )

    return {
      groups,
    }
  },
}
</script>

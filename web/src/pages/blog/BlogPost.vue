<template>
  <div class="relative py-16 bg-white">
    <div class="hidden lg:right-0 lg:block lg:absolute lg:inset-y-0 lg:h-full lg:w-full">
      <div class="relative h-full text-lg max-w-prose mx-auto" aria-hidden="true"></div>
    </div>
    <div class="relative px-4 sm:px-6 lg:px-8">
      <div class="text-lg max-w-prose mx-auto">
        <h1 class="space-y-2">
          <span
            class="block text-base text-center text-yellow-600 font-semibold tracking-wide uppercase"
            >{{ surtitle }}
          </span>
          <span
            class="block text-3xl text-center leading-8 font-extrabold tracking-tight text-gray-900 sm:text-4xl"
          >
            {{ title }}
          </span>
          <span
            class="block text-base text-center text-yellow-600 font-semibold tracking-wide uppercase"
          >
            {{ subtitle }}
          </span>
        </h1>

        <a v-if="author" class="flex-shrink-0 group block mt-8">
          <div class="flex items-center">
            <div>
              <a :href="author.link">
                <img class="inline-block h-9 w-9 rounded-full" :src="author.avatar" alt="" />
              </a>
            </div>
            <div class="ml-3">
              <a
                :href="author.link"
                class="text-sm font-medium text-gray-700 group-hover:text-gray-900 block"
              >
                {{ author.name }}
              </a>
              <a
                :href="author.link"
                class="text-xs font-medium text-gray-500 group-hover:text-gray-700 block"
              >
                {{ date }}
              </a>
              <p v-if="readingTime" class="text-xs font-medium text-gray-500">
                Reading time: {{ readingTime }}
              </p>
            </div>
          </div>
        </a>

        <p class="mt-8 text-xl text-gray-500 leading-8 prose-yellow">
          <slot name="introduction"></slot>
        </p>
      </div>
      <div class="mt-6 prose prose-yellow max-w-3xl text-gray-600 mx-auto">
        <slot></slot>
      </div>

      <div class="mt-12 prose prose-yellow prose-lg text-gray-500 mx-auto">
        <h4>More from the blog</h4>
        <ul class="mt-4 space-y-4">
          <li v-for="item in recentPosts" :key="item.name" class="text-base truncate">
            <router-link
              :to="{ name: item.name }"
              class="font-medium text-gray-900 hover:text-gray-700"
            >
              {{ item.meta.blog.title }}
            </router-link>
          </li>
        </ul>
      </div>

      <DiveIn v-if="diveInBanner" />
    </div>
  </div>
</template>

<script>
import { useRouter } from 'vue-router'
import { useHead } from '@vueuse/head'
import DiveIn from '../../components/index/DiveIn.vue'

export default {
  components: { DiveIn },
  props: {
    title: {
      type: String,
      required: true,
    },
    surtitle: {
      type: String,
      required: false,
    },
    subtitle: {
      type: String,
      required: false,
    },
    date: {
      type: String,
      required: true,
    },
    description: {
      type: String,
      required: true,
    },
    author: {
      type: Object,
      required: false,
    },
    readingTime: {
      type: String,
      required: false,
    },
    diveInBanner: {
      type: Boolean,
      default: true,
      required: false,
    },
    image: {
      type: String,
      default: null,
      required: false,
    },
  },
  setup(props) {
    useHead({
      title: props.title + ' | from the Sturdy Blog',
      meta: [
        { property: 'og:title', content: props.title + ' | from the Sturdy Blog' },
        { property: 'description', content: props.description },
        { property: 'og:description', content: props.description },
        ...(props.image ? [{ property: 'og:image', content: props.image }] : []),
      ],
    })

    let routes = useRouter().getRoutes()
    return {
      recentPosts: routes.filter((r) => r.meta.blog),
    }
  },
}
</script>

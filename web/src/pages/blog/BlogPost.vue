<template>
  <Documentation>
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
              <span class="text-xs font-medium text-gray-500 group-hover:text-gray-700 block">
                {{ date }}
              </span>
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
              {{ item.meta.blogTitle }}
            </router-link>
          </li>
        </ul>
      </div>

      <DiveIn v-if="diveInBanner" />
    </div>
  </Documentation>
</template>

<script lang="ts" setup>
import { useRouter } from 'vue-router'
import { useHead } from '@vueuse/head'
import DiveIn from '../../components/index/DiveIn.vue'
import Documentation from '../../layouts/Documentation.vue'
import { defineProps, withDefaults } from 'vue'

interface Author {
  name: string
  avatar: string
  link: string
}

interface Props {
  title: string
  surtitle?: string
  subtitle?: string
  date: string
  description: string
  author?: Author
  readingTime?: string
  diveInBanner?: boolean
  image?: string
}

const props = withDefaults(defineProps<Props>(), {
  diveInBanner: true,
  surtitle: undefined,
  subtitle: undefined,
  readingTime: undefined,
  image: undefined,
  author: undefined,
})

useHead({
  title: props.title + ' | from the Sturdy Blog',
  meta: [
    { property: 'og:title', content: props.title + ' | from the Sturdy Blog' },
    { property: 'description', content: props.description },
    { property: 'og:description', content: props.description },
    { name: 'description', content: props.description },
    ...(props.image ? [{ property: 'og:image', content: props.image }] : []),
  ],
})

let routes = useRouter().getRoutes()
let recentPosts = routes.filter((r) => r.meta.blog)
</script>

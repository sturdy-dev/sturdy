<template>
  <StaticPage
    category="blog"
    title="The latest posts from Sturdy"
    metadescription="Blog posts from Sturdy"
    image="https://images.unsplash.com/photo-1611159063981-b8c8c4301869?ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&ixlib=rb-1.2.1&auto=format&fit=crop&w=668&q=80"
  >
    <div
      class="mt-5 prose prose-yellow text-gray-500 mx-auto lg:max-w-none lg:row-start-1 lg:col-start-1"
    >
      <ul>
        <li v-for="link in links" :key="link.name">
          <router-link :to="{ name: link.name }" class="!no-underline">
            {{ link.title }}
          </router-link>
        </li>
      </ul>
    </div>
  </StaticPage>
</template>

<script lang="ts" setup>
import StaticPage from '../../layouts/StaticPage.vue'
import { useRouter } from 'vue-router'

type nameTitle = { name: string; title: string }

let routes = useRouter().getRoutes()
let links = routes
  .filter((r) => r.meta.blog && r.meta.blogTitle && r.name)
  .map((r): nameTitle => {
    return {
      name: r.name as string,
      title: r.meta.blogTitle as string,
    }
  })
  .filter(Boolean)
</script>

<template>
  <div class="bg-white py-8">
    <div class="relative max-w-7xl mx-auto py-16 px-4 sm:px-6 lg:px-8 overflow-x-hidden">
      <div
        v-if="!wide"
        class="hidden lg:block bg-gray-50 absolute top-0 bottom-0 left-3/4 w-screen"
      />
      <div class="mx-auto text-base max-w-prose lg:grid lg:grid-cols-2 lg:gap-8 lg:max-w-none">
        <div>
          <h2 class="text-base text-yellow-400 font-semibold tracking-wide uppercase">
            {{ category }}
          </h2>
          <h3
            class="mt-2 text-3xl leading-8 font-extrabold tracking-tight text-gray-900 sm:text-4xl"
          >
            {{ title }}
          </h3>
          <h4 v-if="subtitle" class="mt-4 text-2xl leading-8 text-gray-700 sm:text-xl">
            {{ subtitle }}
          </h4>
        </div>
      </div>
      <div class="mt-8 lg:grid lg:gap-8" :class="grid">
        <div v-if="!wide" class="relative lg:row-start-1" :class="imageColClasses">
          <svg
            class="hidden lg:block absolute top-0 right-0 -mt-20 -mr-20"
            width="404"
            height="384"
            fill="none"
            viewBox="0 0 404 384"
            aria-hidden="true"
          >
            <defs>
              <pattern
                id="de316486-4a29-4312-bdfc-fbce2132a2c1"
                x="0"
                y="0"
                width="20"
                height="20"
                patternUnits="userSpaceOnUse"
              >
                <rect x="0" y="0" width="4" height="4" class="text-gray-200" fill="currentColor" />
              </pattern>
            </defs>
            <rect width="404" height="384" fill="url(#de316486-4a29-4312-bdfc-fbce2132a2c1)" />
          </svg>
          <div class="relative text-base mx-auto max-w-prose lg:max-w-none">
            <figure>
              <div class="aspect-w-12 aspect-h-7 lg:aspect-none">
                <img
                  class="rounded-lg shadow-lg object-cover object-center"
                  :src="image"
                  width="1184"
                  height="1376"
                />
              </div>
            </figure>
          </div>
        </div>
        <div class="mt-8 lg:mt-0 lg:col-start-1" :class="textColClasses">
          <slot></slot>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { useHead } from '@vueuse/head'

export default {
  props: {
    title: String,
    subtitle: String,
    category: String,
    narrow: Boolean,
    wide: Boolean,
    metadescription: String,
    image: {
      type: String,
      default() {
        return 'https://images.unsplash.com/photo-1553877522-43269d4ea984?ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&ixlib=rb-1.2.1&auto=format&fit=crop&w=1950&q=80'
      },
    },
  },
  setup(props) {
    useHead({
      title: props.title + ' | Sturdy',
      meta: [
        { property: 'og:title', content: props.title + ' | Sturdy' },
        { property: 'description', content: props.metadescription },
        { property: 'og:description', content: props.metadescription },
      ],
    })
  },
  computed: {
    grid() {
      if (this.narrow) {
        return 'lg:grid-cols-2'
      }
      if (this.wide) {
        return 'lg:grid-cols-1'
      }
      return 'lg:grid-cols-3'
    },
    imageColClasses() {
      if (this.narrow) {
        return 'lg:col-start-2'
      }
      if (this.wide) {
        return '' // no image
      }
      return 'lg:col-start-3'
    },
    textColClasses() {
      // :class="[narrow ? 'lg:col-end-2' : 'lg:col-end-3']"
      if (this.narrow) {
        return 'lg:col-end-2'
      }
      if (this.wide) {
        return 'lg:col-end-1'
      }
      return 'lg:col-end-3'
    },
  },
}
</script>

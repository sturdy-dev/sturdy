<template>
  <div class="relative">
    <div>
      <h2 class="text-3xl font-bold sm:text-4xl">Recent publications</h2>
      <!--<p class="mt-3 text-xl text-gray-500 sm:mt-4">
          Nullam risus blandit ac aliquam justo ipsum. Quam mauris volutpat massa dictumst amet.
          Sapien tortor lacus arcu.
        </p>-->
    </div>
    <div class="grid gap-16 pt-12 lg:grid-cols-3 lg:gap-x-5 lg:gap-y-12">
      <div v-for="post in posts" :key="post.title">
        <div>
          <span class="inline-block">
            <span
              :class="[
                post.category.color,
                'inline-flex items-center px-3 py-0.5 rounded-full text-sm font-medium',
              ]"
            >
              {{ post.category.name }}
            </span>
          </span>
        </div>
        <a :href="post.href" class="block mt-4">
          <p class="text-xl font-semibold text-gray-900">
            {{ post.title }}
          </p>
          <p class="mt-3 text-base text-gray-500">
            {{ post.description }}
          </p>
        </a>
        <div class="mt-6 flex items-center">
          <div class="flex-shrink-0">
            <span>
              <span class="sr-only">{{ post.author.name }}</span>
              <img
                :class="[post.author.rounded ? ' w-10 rounded-full' : '', 'h-10']"
                :src="post.author.imageUrl"
                alt=""
              />
            </span>
          </div>
          <div class="ml-3">
            <p class="text-sm font-medium text-gray-900">
              <span>
                {{ post.author.name }}
              </span>
            </p>
            <div class="flex space-x-1 text-sm text-gray-500">
              <time :datetime="post.datetime">
                {{ post.date }}
              </time>
              <span aria-hidden="true"> &middot; </span>
              <span> {{ post.readingTime }} read </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import dagensIndustriLogo from '../../assets/dagens-industri.png'

const categoryNews = {
  name: 'News',
  href: '#',
  color: 'bg-green-100 text-green-800',
}

const categoryBlog = { name: 'Blog', href: '#', color: 'bg-blue-100 text-blue-800' }

const authorGustav = {
  name: 'Gustav Westling',
  href: '#',
  imageUrl: 'https://getsturdy.com/assets/gustav.e261d7c0.jpeg',
  rounded: true,
}

const posts = [
  {
    title: 'Swedish software startup challenges Microsoft - funded by heavyweights',
    category: categoryNews,
    href: 'https://www.di.se/digital/svensk-programmertjanst-utmanar-microsoft-tar-in-tungviktare/',
    description:
      'In Swedish: Utveckling kräver alltmer samarbete, överblick och möjligheter att kunna backa stegvis i koden. På området finns i dag Microsofts jätte Github som nu får en konkurrent i den svenska uppstickaren Sturdy – med hjälp av en av Githubs medgrundare.',
    date: 'Jul 16, 2021',
    datetime: '2021-07-16',
    author: {
      name: 'Dagens Industri',
      href: 'https://www.di.se/digital/svensk-programmertjanst-utmanar-microsoft-tar-in-tungviktare/',
      imageUrl: dagensIndustriLogo,
    },
    readingTime: '2 min',
  },
  {
    title: 'Humane Code Review',
    category: categoryBlog,
    href: 'https://sturdy.substack.com/p/005-humane-code-review',
    description:
      'We’ve been thinking a lot about what it means to work on a team lately. Most often, developers are working on the same codebase, but are not really working together...',
    date: 'Jun 10, 2021',
    datetime: '2021-06-10',
    author: authorGustav,
    readingTime: '6 min',
  },
  {
    title: 'Importing from Git!',
    category: categoryBlog,
    href: 'https://sturdy.substack.com/p/004-importing-from-git',
    description:
      'First up is comments! You can now write comments on code, both code that you have written as well as comments from others. We’re starting small, and intend to build more interesting workflows for commenting and collaborating on code, so this is just the first step...',
    date: 'May 6, 2021',
    datetime: '2021-05-06',
    author: authorGustav,
    readingTime: '5 min',
  },
]

export default {
  setup() {
    return {
      posts,
    }
  },
}
</script>

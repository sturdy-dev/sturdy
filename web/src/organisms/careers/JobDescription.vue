<template>
  <JobHeader :title="title" />

  <div class="max-w-7xl mx-auto py-16 px-4 sm:px-6 lg:py-20 lg:px-8">
    <div class="flex flex-col-reverse lg:flex-row">
      <div class="lg:flex-1">
        <div class="prose max-w-full p-8">
          <slot></slot>

          <h2>Benefits</h2>

          <ul>
            <li>Work out of our great office in Stockholm, Sweden</li>
            <li>Flexible working hours – a good work/life balance is really important to us</li>
            <li>Competitive compensation, stock options, and pension</li>
            <li>30 days of paid vacation / year</li>
          </ul>

          <h2>What we use</h2>

          <ul>
            <li>
              Sturdy is built in the public (On
              <a href="https://getsturdy.com/sturdy-zyTDsnY">Sturdy</a>, with a mirror to
              <a href="https://github.com/sturdy-dev/sturdy">GitHub</a>)
            </li>
            <li>Our app and frontend is built with Vue, urql, TailwindCSS, and Electron</li>
            <li>The backend is built with Go, PostgreSQL, and GraphQL</li>
          </ul>

          <h2>How to apply</h2>

          <p>Does this role sound like a good fit? Email us at gustav@getsturdy.com</p>

          <ul>
            <li>Include the role's title in your subject line.</li>
            <li>Send along links that best showcase the relevant things you've built and done.</li>
          </ul>
        </div>
      </div>
      <div class="bg-gray-50 lg:w-80 p-8 rounded-xl">
        <dl class="space-y-12">
          <div v-for="faq in faqs" :key="faq.question">
            <dt class="text-lg leading-6 font-medium text-gray-900">
              {{ faq.question }}
            </dt>
            <!-- eslint-disable-next-line vue/no-v-html -->
            <dd class="mt-2 text-base text-gray-500" v-html="faq.answer"></dd>
          </div>
        </dl>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import JobHeader from './JobHeader.vue'
import { useHead } from '@vueuse/head'

export default defineComponent({
  components: { JobHeader },
  props: {
    title: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    const faqs = [
      {
        question: 'Role',
        answer: props.title,
      },
      {
        question: 'Location',
        answer:
          'Onsite<br><a class="font-medium text-gray-900" href="https://www.google.se/maps/place/Torsgatan+26,+113+21+Stockholm/@59.3366426,18.0446003,17z/data=!3m1!4b1!4m5!3m4!1s0x465f9d6496acf2ed:0x28b94ce2807a0dd1!8m2!3d59.3366426!4d18.046789">Torsgatan 26 (United Spaces)<br>Stockholm, Sweden</a>',
      },
      {
        question: 'Compensation',
        answer: 'Competitive salary and equity',
      },
      {
        question: 'Benefits',
        answer: '30 vacation days / year. Pension.',
      },
      {
        question: 'Funding',
        answer: '$1.5M Seed Round',
      },
      {
        question: 'Funded',
        answer: '2021',
      },
    ]

    useHead({
      title: props.title + ' | Careers at Sturdy',
    })

    return {
      faqs,
    }
  },
})
</script>

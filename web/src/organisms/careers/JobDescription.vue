<template>
  <JobHeader :title="title" />

  <div class="max-w-7xl mx-auto py-16 px-4 sm:px-6 lg:py-20 lg:px-8">
    <div class="flex flex-col-reverse lg:flex-row">
      <div class="lg:flex-1">
        <div class="prose max-w-full p-8">
          <h2 class="text-3xl font-extrabold text-gray-900">About Sturdy</h2>
          <p>At Sturdy we're replacing Git with real-time collaboration in the cloud!</p>

          <p>
            This is a great opportunity to be one of the first employees at a fast moving startup,
            funded by <a href="https://www.ycombinator.com/companies/sturdy/">Y Combinator</a>!
            We’re building tech for techies, and naturally really take pride in both the software
            that we’re writing and the product that we’re building, and hope that you do too!
          </p>

          <h2>Sturdy is growing</h2>

          <p>
            We're looking for Backend Engineers with Go experience to join our team, together with
            the team, you'll play a crucial role in designing Sturdy, and making us known to the
            world.
          </p>

          <h2>What we do</h2>

          <p>
            We're building the future of how coders like you and me are collaborating on code. With
            Sturdy developers don't have to fight branches, Pull Requests, or crazy command line
            arguments to get things done.
          </p>

          <p>
            We regularly post about tech and product on our
            <router-link :to="{ name: 'blog' }">blog</router-link>.
          </p>

          <h2>About the role</h2>
          <p>
            The backend powers the Sturdy real-time version control system that you see in the app.
            It's built with Go, GraphQL, PostgreSQL, and runs on AWS.
          </p>
          <p>
            You’ll be building and architecting a distributed system to build a kick-ass product
            that always works, and always makes our users delighted.
          </p>

          <h2>Qualifications</h2>

          <ul>
            <li>
              Desire to learn or existing knowledge of <strong>Go</strong> (concurrency, error
              handling, testing)
            </li>
            <li>Experience of building distributed high-performance systems</li>
            <li>
              Knowledge of <strong>Git</strong>
              <em>
                (we're using Sturdy to build Sturdy, but the VCS behind the scenes is largely
                inspired by Git's fundamentals)
              </em>
            </li>
            <li>Experience of SQL-databases (we're using PostgreSQL-compatible AWS Aurora)</li>
          </ul>

          <h2>Our tech stack</h2>

          <div class="grid grid-cols-2">
            <div>
              <h3>Backend</h3>
              <ul>
                <li><strong>Sturdy</strong></li>
                <li>A monolithic <strong>Go</strong> service</li>
                <li><strong>PostgreSQL</strong>-compatible AWS Aurora</li>
                <li><strong>GraphQL</strong></li>
                <li>AWS SQS, SES, EC2 (and more)</li>
                <li>Prometheus, Grafana, PostHog</li>
              </ul>
            </div>

            <div>
              <h3>Frontend</h3>
              <ul>
                <li><strong>Vue 3</strong></li>
                <li>TypeScript</li>
                <li>GraphQL (urql)</li>
                <li>Electron</li>
                <li>Tailwind CSS</li>
              </ul>
            </div>
          </div>

          <p>
            We're not expecting you to have experience of our entire tech-stack. But we're expecting
            you to be willing to learn it, and to help us making it even better!
          </p>

          <h2>Ready to apply?</h2>
          <p>
            Send us an email to <a href="mailto:gustav@getsturdy.com">gustav@getsturdy.com</a>, and
            attach a CV/resumé.
          </p>
        </div>
      </div>
      <div class="bg-gray-50 lg:w-80 p-8 rounded-xl">
        <dl class="space-y-12">
          <div v-for="faq in faqs" :key="faq.question">
            <dt class="text-lg leading-6 font-medium text-gray-900">
              {{ faq.question }}
            </dt>
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

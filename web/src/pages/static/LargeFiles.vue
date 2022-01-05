<template>
  <StaticPage
    title="Large files"
    category="solutions"
    metadescription="Sturdy supports large files"
    image="https://images.unsplash.com/photo-1547194936-28214bd75193?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=3300&q=80"
  >
    <div class="text-base max-w-prose mx-auto lg:max-w-none">
      <p class="text-lg text-gray-500">How Sturdy supports enormous files</p>
    </div>
    <div
      class="mt-5 prose prose-yellow text-gray-500 mx-auto lg:max-w-none lg:row-start-1 lg:col-start-1"
    >
      <p>
        Sturdy stores artifacts and source code side-by-side. Artifacts, such as large images or
        binaries can be stored and be a part of the same change as the related source code. Keeping
        everything in one central location makes the workflow simpler.
      </p>

      <p>
        Support for large files is always available, and you don't need any additional
        configuration.
      </p>

      <h2>Made for scalability</h2>
      <p>
        With Sturdy, connected computers only download the data that they need, and does not force
        you to download the entire history of your project.
      </p>
      <p>Sturdy is optimized for making many changes, in massive codebases.</p>

      <h2>Works out of the box</h2>
      <p>
        Compared to Git with Git LFS, Sturdy does not require you to install any additional software
        or configure software on the client side.
      </p>
      <p>The support for large files in Sturdy works out of the box.</p>
    </div>

    <!-- Feature comparison (up to lg) -->
    <section aria-labelledby="mobile-comparison-heading" class="lg:hidden">
      <h2 id="mobile-comparison-heading" class="sr-only">Feature comparison</h2>

      <div class="max-w-2xl mx-auto py-16 px-4 space-y-16 sm:px-6">
        <div
          v-for="(plan, mobilePlanIndex) in plans"
          :key="mobilePlanIndex"
          class="border-t border-gray-200"
        >
          <div
            :class="[
              plan.featured ? 'border-yellow-400' : 'border-transparent',
              '-mt-px pt-6 border-t-2 sm:w-1/2',
            ]"
          >
            <h3 :class="[plan.featured ? 'text-yellow-400' : 'text-gray-900', 'text-sm font-bold']">
              {{ plan.title }}
            </h3>
          </div>

          <div class="mt-6 relative">
            <!-- Fake card background -->
            <div aria-hidden="true" class="hidden absolute inset-0 pointer-events-none sm:block">
              <div
                :class="[
                  plan.featured ? 'shadow-md' : 'shadow',
                  'absolute right-0 w-1/2 h-full bg-white rounded-lg',
                ]"
              />
            </div>

            <div
              :class="[
                plan.featured
                  ? 'ring-2 ring-yellow-400 shadow-md'
                  : 'ring-1 ring-black ring-opacity-5 shadow',
                'relative py-3 px-4 bg-white rounded-lg sm:p-0 sm:bg-transparent sm:rounded-none sm:ring-0 sm:shadow-none',
              ]"
            >
              <dl class="divide-y divide-gray-200">
                <div
                  v-for="perk in perks"
                  :key="perk.title"
                  class="py-3 flex justify-between sm:grid sm:grid-cols-2"
                >
                  <dt class="text-sm font-medium text-gray-600 sm:pr-4">{{ perk.title }}</dt>
                  <dd class="flex items-center justify-end sm:px-4 sm:justify-center">
                    <span
                      v-if="typeof perk.tiers[mobilePlanIndex].value === 'string'"
                      :class="[
                        perk.tiers[mobilePlanIndex].featured ? 'text-yellow-400' : 'text-gray-900',
                        'text-sm font-medium',
                      ]"
                      >{{ perk.tiers[mobilePlanIndex].value }}</span
                    >
                    <template v-else>
                      <CheckIcon
                        v-if="perk.tiers[mobilePlanIndex].value === true"
                        class="mx-auto h-5 w-5 text-yellow-400"
                        aria-hidden="true"
                      />
                      <XIcon v-else class="mx-auto h-5 w-5 text-gray-400" aria-hidden="true" />
                      <span class="sr-only">{{
                        perk.tiers[mobilePlanIndex].value === true ? 'Yes' : 'No'
                      }}</span>
                    </template>
                    <Pill v-if="perk.tiers[mobilePlanIndex].comingSoon" color="green" class="ml-2"
                      >Coming Soon</Pill
                    >
                  </dd>
                </div>
              </dl>
            </div>

            <!-- Fake card border -->
            <div aria-hidden="true" class="hidden absolute inset-0 pointer-events-none sm:block">
              <div
                :class="[
                  plan.featured ? 'ring-2 ring-yellow-400' : 'ring-1 ring-black ring-opacity-5',
                  'absolute right-0 w-1/2 h-full rounded-lg',
                ]"
              />
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Feature comparison (lg+) -->
    <section aria-labelledby="comparison-heading" class="hidden lg:block">
      <h2 id="comparison-heading" class="sr-only">Feature comparison</h2>

      <div class="max-w-7xl mx-auto py-24">
        <h3 class="mb-10 text-sm font-bold text-gray-900">Sturdy vs Git vs Perforce</h3>

        <div class="w-full border-t border-gray-200 flex items-stretch">
          <div class="-mt-px w-1/4 py-6 pr-4 flex items-end"></div>
          <div
            v-for="(plan, planIdx) in plans"
            :key="plan.title"
            aria-hidden="true"
            :class="[planIdx === plans.length - 1 ? '' : 'pr-4', '-mt-px pl-4 w-1/4']"
          >
            <div
              :class="[
                plan.featured ? 'border-yellow-400' : 'border-transparent',
                'py-6 border-t-2',
              ]"
            >
              <p
                :class="[plan.featured ? 'text-yellow-400' : 'text-gray-900', 'text-sm font-bold']"
              >
                {{ plan.title }}
              </p>
            </div>
          </div>
        </div>

        <div class="mt-6 relative">
          <!-- Fake card backgrounds -->
          <div class="absolute inset-0 flex items-stretch pointer-events-none" aria-hidden="true">
            <div class="w-1/4 pr-4" />
            <div class="w-1/4 px-4">
              <div class="w-full h-full bg-white rounded-lg shadow" />
            </div>
            <div class="w-1/4 px-4">
              <div class="w-full h-full bg-white rounded-lg shadow-md" />
            </div>
            <div class="w-1/4 pl-4">
              <div class="w-full h-full bg-white rounded-lg shadow" />
            </div>
          </div>

          <table class="relative w-full">
            <caption class="sr-only">
              Perk comparison
            </caption>
            <thead>
              <tr class="text-left">
                <th scope="col">
                  <span class="sr-only">Perk</span>
                </th>
                <th v-for="plan in plans" :key="plan.title" scope="col">
                  <span class="sr-only">{{ plan.title }} plan</span>
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-for="perk in perks" :key="perk.title">
                <th scope="row" class="w-1/4 py-3 pr-4 text-left text-sm font-medium text-gray-600">
                  {{ perk.title }}
                </th>
                <td
                  v-for="tier in perk.tiers"
                  :key="tier.title"
                  :class="[' px-4 relative w-1/4 py-0 text-center']"
                >
                  <span class="relative w-full h-full py-3">
                    <span
                      v-if="typeof tier.value === 'string'"
                      :class="[
                        tier.featured ? 'text-yellow-400' : 'text-gray-900',
                        'text-sm font-medium',
                      ]"
                    >
                      {{ tier.value }}
                    </span>
                    <template v-else>
                      <CheckIcon
                        v-if="tier.value === true"
                        class="mx-auto h-5 w-5 text-yellow-400"
                        aria-hidden="true"
                      />
                      <XIcon v-else class="mx-auto h-5 w-5 text-gray-400" aria-hidden="true" />
                      <span class="sr-only">{{ tier.value === true ? 'Yes' : 'No' }}</span>
                    </template>
                    <Pill v-if="tier.comingSoon" color="green" class="mb-1">Coming Soon</Pill>
                  </span>
                </td>
              </tr>
            </tbody>
          </table>

          <!-- Fake card borders -->
          <div class="absolute inset-0 flex items-stretch pointer-events-none" aria-hidden="true">
            <div class="w-1/4 pr-4" />
            <div class="w-1/4 px-4">
              <div class="w-full h-full rounded-lg ring-2 ring-yellow-400 ring-opacity-100" />
            </div>
            <div class="w-1/4 px-4">
              <div class="w-full h-full rounded-lg ring-1 ring-black ring-opacity-5" />
            </div>
            <div class="w-1/4 pl-4">
              <div class="w-full h-full rounded-lg ring-1 ring-black ring-opacity-5" />
            </div>
          </div>
        </div>
      </div>
    </section>
  </StaticPage>
</template>

<script>
import StaticPage from './StaticPage.vue'
import { CheckIcon, XIcon } from '@heroicons/vue/outline'
import Pill from '../../components/shared/Pill.vue'

const plans = [
  {
    title: 'Sturdy',
    featured: true,
  },
  {
    title: 'Git',
    featured: false,
  },
  {
    title: 'Perforce',
    featured: false,
  },
]

const perks = [
  {
    title: 'Large files (> 1GB)',
    tiers: [
      { title: 'sturdy', featured: true, value: true },
      { title: 'git', value: 'With LFS' },
      { title: 'perforce', value: true },
    ],
  },
  {
    title: 'Real time collaboration',
    tiers: [
      { title: 'sturdy', featured: true, value: true },
      { title: 'git', value: false },
      { title: 'perforce', value: false },
    ],
  },
  {
    title: 'Modern code review',
    tiers: [
      { title: 'sturdy', featured: true, value: true },
      { title: 'git', value: 'Available' },
      { title: 'perforce', value: false },
    ],
  },
  {
    title: 'Centralized',
    tiers: [
      { title: 'sturdy', featured: true, value: true },
      { title: 'git', value: false },
      { title: 'perforce', value: true },
    ],
  },
  {
    title: 'Locking',
    tiers: [
      { title: 'sturdy', featured: true, value: 'Automatically', comingSoon: true },
      { title: 'git', value: false },
      { title: 'perforce', value: 'Yes, manually' },
    ],
  },
]

export default {
  components: { StaticPage, CheckIcon, XIcon, Pill },
  setup() {
    return {
      plans,
      perks,
    }
  },
}
</script>

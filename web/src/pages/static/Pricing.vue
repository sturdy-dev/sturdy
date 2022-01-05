<!-- This example requires Tailwind CSS v2.0+ -->
<template>
  <div class="bg-gray-50">
    <div class="relative bg-gradient-to-b from-gray-800">
      <div
        class="relative max-w-2xl mx-auto pt-16 px-4 text-center sm:pt-32 sm:px-6 lg:max-w-7xl lg:px-8"
      >
        <h1 class="text-4xl font-extrabold tracking-tight text-white sm:text-6xl">
          <span class="block lg:inline">Simple pricing,&nbsp;</span>
          <span class="block lg:inline">no commitment.</span>
        </h1>
        <p class="mt-4 text-xl text-white">
          Everything you need, nothing you don't. Pick a plan that best suits your business.
        </p>
      </div>

      <h2 class="sr-only">Plans</h2>

      <!-- Toggle -->
      <div class="relative mt-12 flex justify-center sm:mt-16">
        <div class="bg-yellow-400 p-0.5 rounded-lg flex">
          <button
            type="button"
            :class="[
              monthlyBilling
                ? 'bg-white text-yellow-800 hover:bg-yellow-50'
                : 'text-black hover:bg-yellow-500',
            ]"
            class="relative py-2 px-6 border-yellow-500 rounded-md shadow-sm text-sm font-medium whitespace-nowrap focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-yellow-500 focus:ring-white focus:z-10"
            @click="monthlyBilling = true"
          >
            Monthly billing
          </button>
          <button
            type="button"
            :class="[
              !monthlyBilling
                ? 'bg-white text-yellow-800 hover:bg-yellow-50'
                : 'text-black hover:bg-yellow-500',
            ]"
            class="ml-0.5 relative py-2 px-6 border border-transparent rounded-md text-sm font-medium whitespace-nowrap focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-yellow-500 focus:ring-white focus:z-10"
            @click="monthlyBilling = false"
          >
            Yearly billing
          </button>
        </div>
      </div>

      <!-- Cards -->
      <div
        class="relative mt-8 max-w-2xl mx-auto px-4 pb-8 sm:mt-12 sm:px-6 lg:max-w-7xl lg:px-8 lg:pb-0"
      >
        <!-- Decorative background -->
        <div
          aria-hidden="true"
          class="hidden absolute top-4 bottom-6 left-8 right-8 inset-0 bg-gray-50 rounded-lg lg:block"
        />

        <div class="relative space-y-6 lg:space-y-0 lg:grid lg:grid-cols-3">
          <div
            v-for="plan in plans"
            :key="plan.title"
            :class="[
              plan.featured
                ? 'bg-white ring-2 ring-yellow-500 shadow-md'
                : 'bg-white lg:bg-transparent',
              'pt-6 px-6 pb-3 rounded-lg lg:px-8 lg:pt-12',
            ]"
          >
            <div>
              <h3
                :class="[
                  plan.featured ? 'text-yellow-400' : 'text-black',
                  'text-sm font-semibold uppercase tracking-wide',
                ]"
              >
                {{ plan.title }}
              </h3>
              <div
                class="flex flex-col items-start sm:flex-row sm:items-center sm:justify-between lg:flex-col lg:items-start"
              >
                <div v-if="plan.contactUsPricing" class="mt-3 flex items-center">
                  <p
                    :class="[
                      plan.featured ? 'text-yellow-400' : 'text-black',
                      'leading-10 text-2xl font-extrabold ',
                    ]"
                  >
                    Contact us
                  </p>
                </div>
                <div v-else class="mt-3 flex items-center">
                  <p
                    :class="[
                      plan.featured ? 'text-yellow-400' : 'text-black',
                      'text-4xl font-extrabold tracking-tight',
                    ]"
                  >
                    ${{ monthlyBilling ? plan.priceMonthly : plan.priceYearly }}
                  </p>
                  <div class="ml-4">
                    <p :class="[plan.featured ? 'text-gray-700' : 'text-black', 'text-sm']">
                      per user / {{ monthlyBilling ? 'month' : 'year' }}
                    </p>
                    <p
                      v-if="plan.showYearlyPricing && monthlyBilling"
                      :class="[plan.featured ? 'text-gray-500' : 'text-black', 'text-sm']"
                    >
                      Billed yearly (${{ plan.priceYearly }})
                    </p>
                  </div>
                </div>
                <a
                  :href="plan.href"
                  :class="[
                    plan.featured
                      ? 'bg-yellow-400 text-white hover:bg-yellow-500'
                      : 'bg-black text-white hover:bg-gray-900',
                    'mt-6 w-full inline-block py-2 px-8 border border-transparent rounded-md shadow-sm text-center text-sm font-medium sm:mt-0 sm:w-auto lg:mt-6 lg:w-full',
                  ]"
                  >{{ plan.callToAction }}</a
                >
              </div>
            </div>
            <h4 class="sr-only">Features</h4>
            <ul
              :class="[
                plan.featured
                  ? 'border-gray-200 divide-gray-200'
                  : 'border-yellow-500 divide-black divide-opacity-75',
                'mt-7 border-t divide-y lg:border-t-0',
              ]"
            >
              <li
                v-for="mainFeature in plan.mainFeatures"
                :key="mainFeature.id"
                class="py-3 flex items-center"
              >
                <CheckIcon
                  :class="[
                    plan.featured ? 'text-yellow-500' : 'text-black',
                    'w-5 h-5 flex-shrink-0',
                  ]"
                  aria-hidden="true"
                />
                <span
                  :class="[
                    plan.featured ? 'text-gray-600' : 'text-black',
                    'ml-3 text-sm font-medium',
                  ]"
                  >{{ mainFeature.value }}</span
                >
              </li>
            </ul>
          </div>
        </div>
      </div>
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
            <p class="mt-2 text-sm text-gray-500">{{ plan.description }}</p>
          </div>
          <h4 class="mt-10 text-sm font-bold text-gray-900">Catered for business</h4>

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
                  v-for="feature in features"
                  :key="feature.title"
                  class="py-3 flex items-center justify-between sm:grid sm:grid-cols-2"
                >
                  <dt class="pr-4 text-sm font-medium text-gray-600">{{ feature.title }}</dt>
                  <dd class="flex items-center justify-end sm:px-4 sm:justify-center">
                    <span
                      v-if="typeof feature.tiers[mobilePlanIndex].value === 'string'"
                      :class="[
                        feature.tiers[mobilePlanIndex].featured
                          ? 'text-yellow-400'
                          : 'text-gray-900',
                        'text-sm font-medium',
                      ]"
                      >{{ feature.tiers[mobilePlanIndex].value }}</span
                    >
                    <template v-else>
                      <CheckIcon
                        v-if="feature.tiers[mobilePlanIndex].value === true"
                        class="mx-auto h-5 w-5 text-yellow-400"
                        aria-hidden="true"
                      />
                      <XIcon v-else class="mx-auto h-5 w-5 text-gray-400" aria-hidden="true" />
                      <span class="sr-only">{{
                        feature.tiers[mobilePlanIndex].value === true ? 'Yes' : 'No'
                      }}</span>
                    </template>
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

          <h4 class="mt-10 text-sm font-bold text-gray-900">Support and other perks</h4>

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
                  <dd class="text-center sm:px-4">
                    <CheckIcon
                      v-if="perk.tiers[mobilePlanIndex].value === true"
                      class="mx-auto h-5 w-5 text-yellow-400"
                      aria-hidden="true"
                    />
                    <XIcon v-else class="mx-auto h-5 w-5 text-gray-400" aria-hidden="true" />
                    <span class="sr-only">{{
                      perk.tiers[mobilePlanIndex].value === true ? 'Yes' : 'No'
                    }}</span>
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

      <div class="max-w-7xl mx-auto py-24 px-8">
        <div class="w-full border-t border-gray-200 flex items-stretch">
          <div class="-mt-px w-1/4 py-6 pr-4 flex items-end">
            <h3 class="mt-auto text-sm font-bold text-gray-900">Catered for teams</h3>
          </div>
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
              <p class="mt-2 text-sm text-gray-500">{{ plan.description }}</p>
            </div>
          </div>
        </div>

        <div class="relative">
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
              Business feature comparison
            </caption>
            <thead>
              <tr class="text-left">
                <th scope="col">
                  <span class="sr-only">Feature</span>
                </th>
                <th v-for="plan in plans" :key="plan.title" scope="col">
                  <span class="sr-only">{{ plan.title }} plan</span>
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-for="feature in features" :key="feature.title">
                <th scope="row" class="w-1/4 py-3 pr-4 text-left text-sm font-medium text-gray-600">
                  {{ feature.title }}
                </th>
                <td
                  v-for="tier in feature.tiers"
                  :key="tier.title"
                  :class="['px-4 relative w-1/4 py-0 text-center']"
                >
                  <span class="relative w-full h-full py-3">
                    <span
                      v-if="typeof tier.value === 'string'"
                      :class="[
                        tier.featured ? 'text-yellow-400' : 'text-gray-900',
                        'text-sm font-medium',
                      ]"
                      >{{ tier.value }}</span
                    >
                    <template v-else>
                      <CheckIcon
                        v-if="tier.value === true"
                        class="mx-auto h-5 w-5 text-yellow-400"
                        aria-hidden="true"
                      />
                      <XIcon v-else class="mx-auto h-5 w-5 text-gray-400" aria-hidden="true" />
                      <span class="sr-only">{{ tier.value === true ? 'Yes' : 'No' }}</span>
                    </template>
                  </span>
                </td>
              </tr>
            </tbody>
          </table>

          <!-- Fake card borders -->
          <div class="absolute inset-0 flex items-stretch pointer-events-none" aria-hidden="true">
            <div class="w-1/4 pr-4" />
            <div class="w-1/4 px-4">
              <div class="w-full h-full rounded-lg ring-1 ring-black ring-opacity-5" />
            </div>
            <div class="w-1/4 px-4">
              <div class="w-full h-full rounded-lg ring-2 ring-yellow-400 ring-opacity-100" />
            </div>
            <div class="w-1/4 pl-4">
              <div class="w-full h-full rounded-lg ring-1 ring-black ring-opacity-5" />
            </div>
          </div>
        </div>

        <h3 class="mt-10 text-sm font-bold text-gray-900">Support and other perks</h3>

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
                    <CheckIcon
                      v-if="tier.value === true"
                      class="mx-auto h-5 w-5 text-yellow-400"
                      aria-hidden="true"
                    />
                    <XIcon v-else class="mx-auto h-5 w-5 text-gray-400" aria-hidden="true" />
                    <span class="sr-only">{{ tier.value === true ? 'Yes' : 'No' }}</span>
                  </span>
                </td>
              </tr>
            </tbody>
          </table>

          <!-- Fake card borders -->
          <div class="absolute inset-0 flex items-stretch pointer-events-none" aria-hidden="true">
            <div class="w-1/4 pr-4" />
            <div class="w-1/4 px-4">
              <div class="w-full h-full rounded-lg ring-1 ring-black ring-opacity-5" />
            </div>
            <div class="w-1/4 px-4">
              <div class="w-full h-full rounded-lg ring-2 ring-yellow-400 ring-opacity-100" />
            </div>
            <div class="w-1/4 pl-4">
              <div class="w-full h-full rounded-lg ring-1 ring-black ring-opacity-5" />
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
import { CheckIcon, XIcon } from '@heroicons/vue/solid'
import { useHead } from '@vueuse/head'

const plans = [
  {
    title: 'Free',
    featured: false,
    description: 'All the basics, get started using Sturdy today! No strings attached.',
    priceMonthly: 0,
    priceYearly: 0,
    showYearlyPricing: false,
    callToAction: 'Sign Up',
    href: '/signup',
    // callToAction: 'Coming soon',
    // href: '/',
    mainFeatures: [
      { id: 1, value: 'Unlimited Codebases' },
      { id: 2, value: 'Codebases with up to 3 collaborators' },
      { id: 3, value: 'Live feedback and suggestions' },
    ],
  },
  {
    title: 'Startup',
    featured: true,
    description: 'Everything you need for your team to be productive.',
    showYearlyPricing: true,
    priceMonthly: 12,
    priceYearly: 120,
    callToAction: 'Buy Startup',
    href: '/signup',
    // callToAction: 'Coming soon',
    // href: '/',
    mainFeatures: [
      { id: 1, value: 'Unlimited Codebases' },
      { id: 2, value: 'Unlimited Collaborators' },
      { id: 3, value: 'Live feedback and suggestions' },
      { id: 4, value: 'Continuous Integration' },
      // { id: 4, value: 'Codebases with up to 3 collaborators' },
      // { id: 5, value: 'Advanced invoicing' },
      // { id: 6, value: 'Easy to use accounting' },
      // { id: 7, value: 'Mutli-accounts' },
      // { id: 8, value: 'Tax planning toolkit' },
      // { id: 9, value: 'VAT & VATMOSS filing' },
      // { id: 10, value: 'Free bank transfers' },
    ],
  },
  {
    title: 'Enterprise',
    featured: false,
    description: 'Convenient features to take your organization to the next level.',
    contactUsPricing: true,
    callToAction: 'Contact Sales',
    href: 'mailto:sales@getsturdy.com',
    mainFeatures: [
      { id: 1, value: 'SAML sign-on' },
      { id: 2, value: 'Premium Support' },
      { id: 3, value: 'Advanced Audit' },
      // { id: 3, value: 'Mutli-accounts' },
      // { id: 4, value: 'Tax planning toolkit' },
    ],
  },
]
const features = [
  {
    title: 'Codebases',
    tiers: [
      { title: 'free', value: 'Unlimited codebases' },
      { title: 'startup', featured: true, value: 'Unlimited codebases' },
      { title: 'enterprise', value: 'Unlimited codebases' },
    ],
  },
  {
    title: 'Collaborators',
    tiers: [
      { title: 'free', value: 'Up to 3 per codebase' },
      { title: 'startup', featured: true, value: 'Unlimited collaborators' },
      { title: 'enterprise', value: 'Unlimited collaborators' },
    ],
  },
  {
    title: 'Private codebases',
    tiers: [
      { title: 'free', value: true },
      { title: 'startup', featured: true, value: true },
      { title: 'enterprise', value: true },
    ],
  },
  {
    title: 'Code review',
    tiers: [
      { title: 'free', value: true },
      { title: 'startup', featured: true, value: true },
      { title: 'enterprise', value: true },
    ],
  },
  {
    title: 'Live feedback',
    tiers: [
      { title: 'free', value: true },
      { title: 'startup', featured: true, value: true },
      { title: 'enterprise', value: true },
    ],
  },
  {
    title: 'Real-time suggestions',
    tiers: [
      { title: 'free', value: true },
      { title: 'startup', featured: true, value: true },
      { title: 'enterprise', value: true },
    ],
  },

  {
    title: 'Notifications',
    tiers: [
      { title: 'free', value: true },
      { title: 'startup', featured: true, value: true },
      { title: 'enterprise', value: true },
    ],
  },
  {
    title: 'CI/CD',
    tiers: [
      { title: 'free', value: true },
      { title: 'startup', featured: true, value: true },
      { title: 'enterprise', value: true },
    ],
  },
  {
    title: 'GraphQL API',
    tiers: [
      { title: 'free', value: true },
      { title: 'startup', featured: true, value: true },
      { title: 'enterprise', value: true },
    ],
  },
  {
    title: 'Instant Integration',
    tiers: [
      { title: 'free', value: true },
      { title: 'startup', featured: true, value: true },
      { title: 'enterprise', value: true },
    ],
  },

  /*{
    title: 'Multi-accounts',
    tiers: [
      { title: 'starter', value: '3 accounts' },
      { title: 'popular', featured: true, value: 'Unlimited accounts' },
      { title: 'intermediate', value: '7 accounts' },
    ],
  },
  {
    title: 'Invoicing',
    tiers: [
      { title: 'starter', value: '3 invoices' },
      { title: 'popular', featured: true, value: 'Unlimited invoices' },
      { title: 'intermediate', value: '10 invoices' },
    ],
  },
  {
    title: 'Exclusive offers',
    tiers: [
      { title: 'starter', value: false },
      { title: 'popular', featured: true, value: true },
      { title: 'intermediate', value: true },
    ],
  },
  {
    title: '6 months free advisor',
    tiers: [
      { title: 'starter', value: false },
      { title: 'popular', featured: true, value: true },
      { title: 'intermediate', value: true },
    ],
  },
  {
    title: 'Mobile and web access',
    tiers: [
      { title: 'starter', value: false },
      { title: 'popular', featured: true, value: true },
      { title: 'intermediate', value: false },
    ],
  },*/
]
const perks = [
  {
    title: 'SAML single sign on',
    tiers: [
      { title: 'free', value: false },
      { title: 'startup', featured: true, value: false },
      { title: 'enterprise', value: true },
    ],
  },
  {
    title: 'Advanced audit trail',
    tiers: [
      { title: 'free', value: false },
      { title: 'startup', featured: true, value: false },
      { title: 'enterprise', value: true },
    ],
  },
  {
    title: 'Community Support',
    tiers: [
      { title: 'free', value: true },
      { title: 'startup', featured: true, value: true },
      { title: 'enterprise', value: true },
    ],
  },
  {
    title: 'Standard Support',
    tiers: [
      { title: 'free', value: false },
      { title: 'startup', featured: true, value: true },
      { title: 'enterprise', value: true },
    ],
  },
  {
    title: 'Premium Support',
    tiers: [
      { title: 'free', value: false },
      { title: 'startup', featured: true, value: false },
      { title: 'enterprise', value: true },
    ],
  },
  /*{
    title: 'Digital receipts',
    tiers: [
      { title: 'starter', value: true },
      { title: 'popular', featured: true, value: true },
      { title: 'intermediate', value: true },
    ],
  },
  {
    title: 'Pots to separate money',
    tiers: [
      { title: 'starter', value: false },
      { title: 'popular', featured: true, value: true },
      { title: 'intermediate', value: true },
    ],
  },
  {
    title: 'Free bank transfers',
    tiers: [
      { title: 'starter', value: false },
      { title: 'popular', featured: true, value: true },
      { title: 'intermediate', value: false },
    ],
  },
  {
    title: 'Business debit card',
    tiers: [
      { title: 'starter', value: false },
      { title: 'popular', featured: true, value: true },
      { title: 'intermediate', value: false },
    ],
  },*/
]

export default {
  components: {
    CheckIcon,
    XIcon,
  },
  setup() {
    useHead({
      title: 'Pricing | Sturdy',
    })
    return {
      plans,
      features,
      perks,
    }
  },
  data: function () {
    return {
      monthlyBilling: true,
    }
  },
}
</script>

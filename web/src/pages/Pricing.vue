<template>
  <DocumentationStickyHeader class="antialiased text-slate-300 bg-[#0c1322]">
    <div class="px-6 relative">
      <div class="absolute inset-0 bottom-20 bg-bottom bg-no-repeat gradient"></div>

      <div class="relative max-w-5xl mx-auto pt-24 sm:pt-28 lg:pt-32">
        <h1
          class="font-extrabold text-4xl sm:text-5xl lg:text-6xl tracking-tight text-center text-slate-50"
        >
          Pricing
        </h1>

        <p class="mt-6 text-lg text-slate-400 text-center max-w-3xl mx-auto">
          Find the version of Sturdy that's best for you.
        </p>

        <div class="mx-auto pt-20">
          <div
            class="flex justify-around text-slate-200 my-12 flex-col lg:flex-row space-y-8 lg:space-y-0 p-4 lg:p-0 lg:space-x-8"
          >
            <PricingTierSummary name="Free">
              <PricingTierSummaryItem>Unlimited codebases and workspaces</PricingTierSummaryItem>
              <PricingTierSummaryItem>Self-hosted, or in the cloud</PricingTierSummaryItem>
              <PricingTierSummaryItem>Up to 10 users</PricingTierSummaryItem>
              <PricingTierSummaryItem>Sturdy for GitHub</PricingTierSummaryItem>
              <PricingTierSummaryItem>CI/CD</PricingTierSummaryItem>
            </PricingTierSummary>

            <PricingTierSummary name="Pro" price="30">
              <PricingTierSummaryItemArrow>Everything from Free</PricingTierSummaryItemArrow>
              <PricingTierSummaryItem>Self-hosted, or in the cloud</PricingTierSummaryItem>
              <PricingTierSummaryItem>Unlimited users</PricingTierSummaryItem>
              <PricingTierSummaryItem>SSO/SAML</PricingTierSummaryItem>
              <PricingTierSummaryItem>Shared support Slack</PricingTierSummaryItem>
              <PricingTierSummaryItem>Onboarding and training</PricingTierSummaryItem>
              <PricingTierSummaryItem>Advanced audit & security</PricingTierSummaryItem>
            </PricingTierSummary>
          </div>
        </div>
      </div>
    </div>

    <div class="block top-0 inset-x-0 mt-12 bg-[#0c1322]">
      <div class="mx-auto py-16 sm:py-24 sm:px-6 lg:px-8 max-w-[100rem]">
        <p class="text-center mb-6 uppercase font-bold">Full comparison</p>

        <!-- xs to lg -->
        <div class="mx-auto space-y-16 lg:hidden">
          <section v-for="(tier, tierIdx) in tiers" :key="tier.name">
            <div class="px-4 mb-8">
              <h2 class="text-lg leading-6 font-medium text-slate-200">{{ tier.name }}</h2>
              <p v-if="tier.priceMonthly === 0" class="mt-4">
                <span class="text-4xl font-extrabold text-slate-200">Free</span>
              </p>
              <p v-else class="mt-4">
                <span class="text-4xl font-extrabold text-slate-200">${{ tier.priceMonthly }}</span>
                {{ ' ' }}
                <span class="text-base font-medium text-slate-400">/mo</span>
              </p>
              <p class="mt-4 text-sm text-slate-400">{{ tier.description }}</p>
              <PricingButton :to="{ name: tier.getStartedRoute }">
                Get started with {{ tier.name }}
              </PricingButton>
            </div>

            <table v-for="section in sections" :key="section.name" class="w-full">
              <caption
                class="bg-slate-800 border-t border-slate-700 py-3 px-4 text-sm font-bold text-slate-200 text-left"
              >
                {{
                  section.name
                }}
              </caption>
              <thead>
                <tr>
                  <th class="sr-only" scope="col">Feature</th>
                  <th class="sr-only" scope="col">Included</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-700">
                <tr
                  v-for="feature in section.features"
                  :key="feature.name"
                  class="border-t border-slate-700"
                >
                  <th class="py-5 px-4 text-sm font-normal text-slate-400 text-left" scope="row">
                    {{ feature.name }}
                  </th>
                  <td class="py-5 pr-4">
                    <span
                      v-if="typeof feature.tiers[tier.name] === 'string'"
                      class="block text-sm text-slate-200 text-right"
                    >
                      {{ feature.tiers[tier.name] }}
                    </span>
                    <template v-else>
                      <CheckIcon
                        v-if="feature.tiers[tier.name] === true"
                        class="ml-auto h-5 w-5 text-green-500"
                        aria-hidden="true"
                      />
                      <MinusIcon v-else class="ml-auto h-5 w-5 text-gray-700" aria-hidden="true" />
                      <span class="sr-only">
                        {{ feature.tiers[tier.name] === true ? 'Yes' : 'No' }}
                      </span>
                    </template>
                  </td>
                </tr>
              </tbody>
            </table>

            <div
              :class="[
                tierIdx < tiers.length - 1 ? 'py-5 border-b' : 'pt-5',
                'border-t border-gray-200 px-4',
              ]"
            >
              <PricingButton :to="{ name: tier.getStartedRoute }">
                Get started with {{ tier.name }}
              </PricingButton>
            </div>
          </section>
        </div>

        <!-- lg+ -->
        <div class="hidden lg:block">
          <table class="w-full h-px table-fixed">
            <caption class="sr-only">
              Pricing plan comparison
            </caption>
            <thead>
              <tr>
                <th></th>
                <th
                  :colspan="tiers.length === 4 ? 3 : 2"
                  class="py-4 px-6 text-md font-light text-slate-400 bg-slate-800/60"
                >
                  Self-hosted options
                </th>
                <th class="py-4 px-6 text-md font-light text-slate-400">Hosted solution</th>
              </tr>
              <tr>
                <th
                  class="py-4 px-6 text-sm font-medium text-slate-200 text-left border-1 border-r border-slate-600"
                  scope="col"
                >
                  <span class="sr-only">Feature by</span>
                  <span>Plans</span>
                </th>
                <th
                  v-for="tier in tiers"
                  :key="tier.name"
                  class="w-1/4 py-4 px-6 text-lg leading-6 font-medium text-slate-200 text-left border-1 border-r border-slate-600"
                  scope="col"
                >
                  {{ tier.name }}
                </th>
              </tr>
            </thead>
            <tbody class="border-t border-slate-800 divide-y divide-slate-800">
              <tr>
                <th
                  class="py-8 px-6 text-sm font-medium text-slate-200 text-left align-top"
                  scope="row"
                ></th>
                <td v-for="tier in tiers" :key="tier.name" class="h-full py-8 px-6 align-top">
                  <div class="relative h-full table w-full">
                    <p v-if="tier.priceMonthly === 0">
                      <span class="text-4xl font-extrabold text-slate-200">Free</span>
                    </p>
                    <p v-else>
                      <span class="text-4xl font-extrabold text-slate-200"
                        >${{ tier.priceMonthly }}</span
                      >
                      {{ ' ' }}
                      <span class="text-base font-medium text-slate-400">/user/mo</span>
                    </p>
                    <p class="mt-4 mb-24 text-sm text-slate-300">{{ tier.description }}</p>
                    <PricingButton :to="{ name: tier.getStartedRoute }" class="absolute bottom-0">
                      Get started with <span class="whitespace-nowrap">{{ tier.name }}</span>
                    </PricingButton>
                  </div>
                </td>
              </tr>
              <template v-for="section in sections" :key="section.name">
                <tr>
                  <th
                    class="bg-slate-800 py-3 pl-6 text-sm font-bold text-slate-200 text-left"
                    :colspan="tiers.length + 1"
                    scope="colgroup"
                  >
                    {{ section.name }}
                  </th>
                </tr>
                <tr v-for="feature in section.features" :key="feature.name">
                  <th class="py-5 px-6 text-sm font-normal text-slate-300 text-left" scope="row">
                    {{ feature.name }}
                  </th>
                  <td v-for="tier in tiers" :key="tier.name" class="py-5 px-6">
                    <span
                      v-if="typeof feature.tiers[tier.name] === 'string'"
                      class="block text-sm text-slate-200"
                    >
                      {{ feature.tiers[tier.name] }}
                    </span>
                    <template v-else>
                      <CheckIcon
                        v-if="feature.tiers[tier.name] === true"
                        class="h-5 w-5 text-green-500"
                        aria-hidden="true"
                      />
                      <MinusIcon v-else class="h-5 w-5 text-gray-700" aria-hidden="true" />
                      <span class="sr-only">
                        {{ feature.tiers[tier.name] === true ? 'Included' : 'Not included' }} in
                        {{ tier.name }}
                      </span>
                    </template>
                  </td>
                </tr>
              </template>
            </tbody>
            <tfoot>
              <tr class="border-t border-gray-200">
                <th class="sr-only" scope="row">Choose your plan</th>
                <td v-for="tier in tiers" :key="tier.name" class="pt-5 px-6">
                  <PricingButton :to="{ name: tier.getStartedRoute }" class="inline-block">
                    Get started with <span class="whitespace-nowrap">{{ tier.name }}</span>
                  </PricingButton>
                </td>
              </tr>
            </tfoot>
          </table>
        </div>
      </div>

      <div class="text-center">
        <h2>Looking for the open-source Sturdy?</h2>
        <p class="mt-2">
          <Button v-if="showOpenSource" color="slate" @click="showOpenSource = false">
            Hide open-source
          </Button>
          <Button v-else color="slate" @click="showOpenSource = true">Show open-source</Button>
        </p>
      </div>
    </div>
  </DocumentationStickyHeader>
</template>

<script lang="ts">
import { CheckIcon, MinusIcon } from '@heroicons/vue/solid'
import { defineComponent, Ref, ref } from 'vue'
import DocumentationStickyHeader from '../layouts/DocumentationStickyHeader.vue'
import Button from '../components/shared/Button.vue'
import PricingTierSummary from '../molecules/pricing/PricingTierSummary.vue'
import PricingTierSummaryItem from '../molecules/pricing/PricingTierSummaryItem.vue'
import PricingTierSummaryItemArrow from '../molecules/pricing/PricingTierSummaryItemArrow.vue'
import PricingButton from '../molecules/pricing/PricingButton.vue'

let showOpenSource = ref(false)
let showTrue = ref(true)

const allTiers = [
  {
    name: 'Open Source',
    getStartedRoute: 'v2DocsSelfHosted',
    priceMonthly: 0,
    description: 'Free and open-source.',
    show: showOpenSource,
  },
  {
    name: 'Free',
    getStartedRoute: 'v2download',
    priceMonthly: 0,
    description: 'Run Sturdy yourself, free forever.',
    show: showTrue,
  },
  {
    name: 'Enterprise',
    getStartedRoute: 'v2download',
    priceMonthly: 30,
    description: 'Advanced features, run anywhere.',
    show: showTrue,
  },
  {
    name: 'Cloud',
    getStartedRoute: 'v2download',
    priceMonthly: 30,
    description: 'Get started for free, with Sturdy in the cloud. No credit card required.',
    show: showTrue,
  },
]

interface Tiers {
  'Open Source': boolean | string
  Free: boolean | string
  Enterprise: boolean | string
  Cloud: boolean | string
}

interface Feature {
  name: string
  tiers: Tiers
  show: Ref<boolean>
}

interface Section {
  name: string
  features: Feature[]
}

const allSections: Section[] = [
  {
    name: 'Overview',
    features: [
      {
        name: 'Benefits',
        tiers: {
          'Open Source': 'Free software, enjoy the core of Sturdy!',
          Free: 'Great for small teams',
          Enterprise: 'Easy migration and incremental migration to Sturdy.',
          Cloud: 'Scales as needed, easy pricing.',
        },
        show: showTrue,
      },

      {
        name: 'Pricing',
        tiers: {
          'Open Source': 'Free',
          Free: 'Free',
          Enterprise: '$30/user (minimum 20 users).',
          Cloud: 'Free (up to 10 users), $30/user (no minimum).',
        },
        show: showTrue,
      },
    ],
  },

  {
    name: 'Instance & deployment',
    features: [
      {
        name: 'Hosting',
        tiers: {
          'Open Source': false,
          Free: false,
          Enterprise: false,
          Cloud: true,
        },
        show: showTrue,
      },

      {
        name: 'All code stays on your infrastructure',
        tiers: {
          'Open Source': true,
          Free: true,
          Enterprise: true,
          Cloud: false,
        },
        show: showTrue,
      },

      {
        name: 'First time setup',
        tiers: {
          'Open Source': 'Compile from source',
          Free: 'Instant',
          Enterprise: '~1 day',
          Cloud: 'Instant',
        },
        show: showTrue,
      },

      {
        name: 'Server management',
        tiers: {
          'Open Source': 'Managed by you',
          Free: 'Managed by you',
          Enterprise: "We'll help you",
          Cloud: 'Managed by us',
        },
        show: showTrue,
      },
    ],
  },

  {
    name: 'Limits',
    features: [
      {
        name: 'Users',
        tiers: {
          'Open Source': 'Unlimited',
          Free: '10',
          Enterprise: 'Unlimited',
          Cloud: 'Unlimited (first 10 for free)',
        },
        show: showTrue,
      },
      {
        name: 'Organizations',
        tiers: {
          'Open Source': '1',
          Enterprise: '1',
          Free: '1',
          Cloud: 'Unlimited',
        },
        show: showTrue,
      },
      {
        name: 'License',
        tiers: {
          'Open Source': 'Apache 2',
          Free: 'Sturdy Enterprise License',
          Enterprise: 'Sturdy Enterprise License',
          Cloud: 'Sturdy Enterprise License',
        },
        show: showOpenSource,
      },
      {
        name: 'Billing',
        tiers: {
          'Open Source': false,
          Free: false,
          Enterprise: 'Annually',
          Cloud: 'Monthly',
        },
        show: showTrue,
      },
    ],
  },

  {
    name: 'Features',
    features: [
      {
        name: 'Live feedback',
        tiers: {
          'Open Source': true,
          Free: true,
          Enterprise: true,
          Cloud: true,
        },
        show: showTrue,
      },

      {
        name: 'Private Codebases',
        tiers: {
          'Open Source': true,
          Free: true,
          Enterprise: true,
          Cloud: true,
        },
        show: showTrue,
      },

      {
        name: 'Public Codebases',
        tiers: {
          'Open Source': true,
          Free: true,
          Enterprise: true,
          Cloud: true,
        },
        show: showTrue,
      },

      {
        name: 'Workspaces',
        tiers: {
          'Open Source': true,
          Free: true,
          Enterprise: true,
          Cloud: true,
        },
        show: showTrue,
      },

      {
        name: 'Suggestions',
        tiers: {
          'Open Source': true,
          Free: true,
          Enterprise: true,
          Cloud: true,
        },
        show: showTrue,
      },

      {
        name: 'Sturdy apps for Mac and Windows',
        tiers: {
          'Open Source': true,
          Free: true,
          Enterprise: true,
          Cloud: true,
        },
        show: showTrue,
      },

      {
        name: 'GraphQL API',
        tiers: {
          'Open Source': true,
          Free: true,
          Enterprise: true,
          Cloud: true,
        },
        show: showTrue,
      },

      {
        name: 'File Access Control Lists',
        tiers: {
          'Open Source': true,
          Free: true,
          Enterprise: true,
          Cloud: true,
        },
        show: showTrue,
      },

      {
        name: 'Sign-up',
        tiers: {
          'Open Source': 'Email and password',
          Free: 'Email and password',
          Enterprise: 'Email and password',
          Cloud: 'Email and magic codes',
        },
        show: showTrue,
      },

      {
        name: 'Sturdy for GitHub',
        tiers: {
          'Open Source': false,
          Free: true,
          Enterprise: true,
          Cloud: true,
        },
        show: showTrue,
      },
      {
        name: 'Buildkite',
        tiers: {
          'Open Source': false,
          Free: true,
          Enterprise: true,
          Cloud: true,
        },
        show: showTrue,
      },
    ],
  },

  {
    name: 'Security',
    features: [
      {
        name: 'ACLs',
        tiers: {
          'Open Source': true,
          Free: true,
          Enterprise: true,
          Cloud: true,
        },
        show: showTrue,
      },
      {
        name: 'Advanced audit and logging',
        tiers: {
          'Open Source': false,
          Free: false,
          Enterprise: 'Contact us',
          Cloud: 'Contact us',
        },
        show: showTrue,
      },
      {
        name: 'SSO/SAML',
        tiers: {
          'Open Source': false,
          Free: false,
          Enterprise: 'Contact us',
          Cloud: 'Contact us',
        },
        show: showTrue,
      },
      {
        name: 'User permissions',
        tiers: {
          'Open Source': false,
          Free: false,
          Enterprise: 'Contact us',
          Cloud: 'Contact us',
        },
        show: showTrue,
      },
    ],
  },

  {
    name: 'Support',
    features: [
      {
        name: 'Community Discord',
        tiers: {
          'Open Source': true,
          Free: true,
          Enterprise: true,
          Cloud: true,
        },
        show: showTrue,
      },

      {
        name: 'Slack (dedicated channel)',
        tiers: {
          'Open Source': false,
          Free: false,
          Enterprise: 'Included',
          Cloud: '$1k/month spend or above',
        },
        show: showTrue,
      },

      {
        name: 'Email',
        tiers: {
          'Open Source': false,
          Free: false,
          Enterprise: 'Included',
          Cloud: 'Included on paid tiers',
        },
        show: showTrue,
      },

      {
        name: 'Training sessions',
        tiers: {
          'Open Source': false,
          Free: false,
          Enterprise: '$1k/month spend or above',
          Cloud: '$1k/month spend or above',
        },
        show: showTrue,
      },

      {
        name: 'Payment via invoicing',
        tiers: {
          'Open Source': false,
          Free: false,
          Enterprise: '$2k/month spend or above',
          Cloud: '$2k/month spend or above',
        },
        show: showTrue,
      },
    ],
  },
]

export default defineComponent({
  components: {
    PricingButton,
    PricingTierSummaryItem,
    PricingTierSummaryItemArrow,
    PricingTierSummary,
    DocumentationStickyHeader,
    CheckIcon,
    MinusIcon,
    Button,
  },
  setup() {
    return {
      showOpenSource,
    }
  },
  computed: {
    tiers() {
      return allTiers.filter((t) => t.show.value)
    },
    sections() {
      return allSections.map((s: Section): Section => {
        return {
          name: s.name,
          features: s.features.filter((f: Feature) => f.show.value),
        }
      })
    },
  },
})
</script>

<style scoped>
.gradient {
  background-image: #0b1120;
  background-image: radial-gradient(at 23% 84%, hsla(223, 49%, 9%, 1) 0, transparent 56%),
    radial-gradient(at 81% 30%, hsla(223, 49%, 9%, 1) 0, transparent 46%),
    radial-gradient(at 87% 34%, hsla(223, 49%, 9%, 1) 0, transparent 57%),
    radial-gradient(at 16% 71%, hsla(223, 49%, 9%, 1) 0, transparent 40%),
    radial-gradient(at 100% 100%, hsla(223, 49%, 9%, 1) 0, transparent 52%),
    radial-gradient(at 51% 63%, hsla(271, 92%, 66%, 1) 0, transparent 49%),
    radial-gradient(at 0% 5%, hsla(223, 49%, 9%, 1) 0, transparent 49%),
    radial-gradient(at 74% 76%, hsla(44, 97%, 57%, 1) 0, transparent 50%),
    radial-gradient(at 45% 37%, hsla(38, 93%, 51%, 1) 0, transparent 50%);
}
</style>

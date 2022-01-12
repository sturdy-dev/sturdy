<template>
  <RadioGroup v-model="selectedTier">
    <RadioGroupLabel class="text-base font-medium text-gray-900"> Select tier </RadioGroupLabel>

    <div class="mt-4 grid grid-cols-1 gap-y-6 sm:grid-cols-3 sm:gap-x-4">
      <RadioGroupOption
        as="template"
        v-for="tier in tiers"
        :key="tier.id"
        :value="tier"
        v-slot="{ checked, active }"
      >
        <div
          :class="[
            checked ? 'border-transparent' : 'border-gray-300',
            active ? 'ring-2 ring-blue-500' : '',
            'relative bg-white border rounded-lg shadow-sm p-4 flex cursor-pointer focus:outline-none',
          ]"
        >
          <div class="flex-1 flex">
            <div class="flex flex-col">
              <RadioGroupLabel as="span" class="block text-sm font-medium text-gray-900">
                {{ tier.title }}
              </RadioGroupLabel>

              <RadioGroupDescription as="span" class="mt-1 flex items-center text-sm text-gray-500">
                <ul class="list-disc list-inside">
                  <li v-for="feat in tier.features" :key="feat">{{ feat }}</li>
                </ul>
              </RadioGroupDescription>

              <RadioGroupDescription as="div" class="mt-6 text-sm font-medium flex items-center">
                <span class="text-3xl text-gray-900">${{ tier.pricePerUser }}</span>
                <span class="text-gray-400">&nbsp;per user / month</span>
              </RadioGroupDescription>

              <RadioGroupDescription
                as="span"
                class="mt-1 flex items-center text-sm text-gray-500"
                v-if="tier.isStartTrial"
              >
                <span>Start 14-day free trial, no credit card required</span>
              </RadioGroupDescription>

              <RadioGroupDescription
                as="span"
                class="mt-1 flex items-center text-sm text-gray-500"
                v-if="tier.minimumUsers"
              >
                <span>Starts at {{ tier.minimumUsers }} users, billed annually.</span>
              </RadioGroupDescription>
            </div>
          </div>
          <component
            :is="tier.checkedIcon"
            :class="[!checked ? 'invisible' : '', 'h-5 w-5 text-blue-600']"
            aria-hidden="true"
          />
          <div
            :class="[
              active ? 'border' : 'border-2',
              checked ? 'border-blue-500' : 'border-transparent',
              'absolute -inset-px rounded-lg pointer-events-none',
            ]"
            aria-hidden="true"
          />
        </div>
      </RadioGroupOption>
    </div>
  </RadioGroup>
</template>

<script>
import { ref } from 'vue'
import {
  RadioGroup,
  RadioGroupDescription,
  RadioGroupLabel,
  RadioGroupOption,
} from '@headlessui/vue'
import { CheckCircleIcon, BriefcaseIcon } from '@heroicons/vue/solid'

const tiers = [
  {
    id: 'free',
    title: 'Free',

    features: ['Unlimited codebases', 'Real-time code review', 'Up to 10 collaborators'],
    pricePerUser: 0,
    checkedIcon: CheckCircleIcon,
  },
  {
    id: 'cloud',
    title: 'Pro',
    features: ['Everything in Free plus...', 'Run on top of GitHub', 'Unlimited Users'],
    pricePerUser: 10,
    isStartTrial: true,
    checkedIcon: CheckCircleIcon,
  },
  {
    id: 'enterprise',
    title: 'Enterprise',
    features: ['Everything in Pro plus...', 'Self-hosted', 'Dedicated Support Slack'],
    pricePerUser: 10,
    minimumUsers: 20,
    checkedIcon: BriefcaseIcon,
  },
]

export default {
  components: {
    RadioGroup,
    RadioGroupDescription,
    RadioGroupLabel,
    RadioGroupOption,
    CheckCircleIcon,
    BriefcaseIcon,
  },
  setup() {
    const selectedTier = ref(tiers[1])

    return {
      tiers,
      selectedTier,
    }
  },
}
</script>

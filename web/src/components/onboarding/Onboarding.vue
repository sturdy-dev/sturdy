<template>
  <transition
    enter-active-class="transition ease-out duration-100 delay-1000"
    enter-from-class="transform opacity-0"
    enter-to-class="transform opacity-100"
    leave-active-class="transition ease-in duration-75"
    leave-from-class="transform opacity-100"
    leave-to-class="transform opacity-0"
  >
    <svg
      v-if="currentStep"
      class="fixed top-0 left-0 hidden md:block"
      height="100%"
      width="100%"
      :viewBox="`0 0 ${screenWidth} ${screenHeight}`"
    >
      <mask id="frame">
        <rect x="0" y="0" :width="screenWidth" :height="screenHeight" fill="white" />
        <rect
          :x="x - 10"
          :y="y - 10"
          :width="width + 20"
          :height="height + 20"
          fill="black"
          rx="15"
          class="transition-all duration-600"
        />
      </mask>
      <rect
        x="0"
        y="0"
        :width="screenWidth"
        :height="screenHeight"
        fill="black"
        opacity="0.4"
        mask="url(#frame)"
      />
      <path :d="bubblePath" fill="#fcfcfc" class="drop-shadow-lg transition-all duration-600" />
      <foreignObject
        class="transition-all duration-600"
        :x="bubbleContentPosition.x"
        :y="bubbleContentPosition.y"
        :height="bubbleHeight + 5"
        :width="bubbleWidth"
        overflow="visible"
      >
        <div class="h-full flex flex-col justify-between">
          <div>
            <div class="flex flex-row items-center gap-1">
              <div class="font-medium text-sm flex-1">
                <RenderVnodes :nodes="currentStep.title()" />
              </div>
              <div class="flex-none text-xs mx-2">{{ steps.progress + 1 }}/{{ steps.length }}</div>
              <button class="flex-none w-5 h-5" @click="cancel()">
                <XIcon />
              </button>
            </div>
            <div class="text-sm">
              <RenderVnodes :nodes="currentStep.description()" />
            </div>
          </div>
          <div class="flex justify-end gap-2">
            <Button :disabled="steps.progress === 0" @click="previous()">Previous </Button>
            <Button :auto-focus="true" @click="next()">
              <template v-if="steps.progress === steps.length - 1">Done </template>
              <template v-else>Next</template>
            </Button>
          </div>
        </div>
      </foreignObject>
    </svg>
  </transition>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { XIcon } from '@heroicons/vue/solid'
import Button from '../shared/Button.vue'
import RenderVnodes from './RenderVnodes'
import { Step } from './Step'
import { StepQueue } from './StepQueue'
import { CombinedError, gql, useClientHandle, useMutation, useSubscription } from '@urql/vue'
import {
  CompletedOnboardingStepsQuery,
  CompletedOnboardingStepsQueryVariables,
  CompletedOnboardingStepSubscription,
  CompletedOnboardingStepSubscriptionVariables,
  OnboardingStepCompletedMutation,
  OnboardingStepCompletedMutationVariables,
} from './__generated__/Onboarding'

const ONBOARDING_STEP_COMPLETED = gql`
  mutation OnboardingStepCompleted($stepID: ID!) {
    completeOnboardingStep(stepID: $stepID) {
      id
    }
  }
`

const COMPLETED_ONBOARDING_STEPS = gql`
  query CompletedOnboardingSteps {
    completedOnboardingSteps {
      id
    }
  }
`

const COMPLETED_ONBOARDING_STEP = gql`
  subscription CompletedOnboardingStep {
    completedOnboardingStep {
      id
    }
  }
`
const errIsUnauthenticated = function (err: undefined | CombinedError): boolean {
  if (!err) return false
  if (!err.graphQLErrors) return false
  return err.graphQLErrors.filter(({ message }) => message === 'UnauthenticatedError').length > 0
}

let steps: StepQueue | undefined

export function useOnboarding() {
  const { executeMutation } = useMutation<
    OnboardingStepCompletedMutation,
    OnboardingStepCompletedMutationVariables
  >(ONBOARDING_STEP_COMPLETED)
  if (steps == null) {
    steps = new StepQueue(async (stepID) => {
      await executeMutation({ stepID })
    })
  }

  const client = useClientHandle()

  function unregisterStep(step: Step) {
    steps?.unregisterStep(step)
  }

  useSubscription<
    CompletedOnboardingStepSubscription,
    void,
    CompletedOnboardingStepSubscriptionVariables
  >(
    {
      query: COMPLETED_ONBOARDING_STEP,
    },
    (prev, data) => {
      if (steps == null) {
        return
      }
      if (steps.currentStep?.id === data.completedOnboardingStep.id) {
        steps.next()
      } else {
        steps.complete(data.completedOnboardingStep.id)
      }
    }
  )

  return {
    steps,
    registerStep(step: Step): boolean {
      const stepRect = calculateChildrenRect(step.highlightedElement)
      if (stepRect.width < 1 || stepRect.y > window.innerHeight) {
        return false
      }
      client.client
        .query<CompletedOnboardingStepsQuery, CompletedOnboardingStepsQueryVariables>(
          COMPLETED_ONBOARDING_STEPS
        )
        .toPromise()
        .then((result) => {
          if (errIsUnauthenticated(result.error)) return

          if (result.data && !result.data.completedOnboardingSteps.some((s) => s.id === step.id)) {
            steps?.registerStep(step)
          }
        })
      return true
    },
    unregisterStep,
  }
}

function calculateChildrenRect(element: HTMLElement): {
  x: number
  y: number
  width: number
  height: number
} {
  let x = Infinity
  let y = Infinity
  let width = 0
  let height = 0
  for (const child of element.children) {
    const rect = child.getBoundingClientRect()
    x = Math.min(x, rect.x)
    y = Math.min(y, rect.y)
    width = Math.max(width, rect.width)
    height = Math.max(height, rect.height)
  }
  return { y, x, width, height }
}

export default defineComponent({
  components: { RenderVnodes, Button, XIcon },
  setup() {
    const { steps } = useOnboarding()
    return { steps, ...steps.refs }
  },
  data() {
    return {
      x: -10,
      y: -10,
      width: 0,
      height: 0,
      screenWidth: window.innerWidth,
      screenHeight: window.innerHeight,
      bubbleWidth: 200,
      bubbleHeight: 250,
    }
  },
  computed: {
    currentStep(): Step | undefined {
      return this.steps.currentStep
    },

    centerTopAligned(): boolean {
      return this.y + this.height + this.bubbleHeight > this.screenHeight
    },

    xDirection(): 'right' | 'left' {
      if (this.x + this.width > this.screenWidth - this.bubbleWidth) {
        return 'left'
      }
      return 'right'
    },

    bubbleContentPosition(): { x: number; y: number } {
      if (this.centerTopAligned) {
        return {
          x: this.x + this.width / 2 - this.bubbleWidth / 2,
          y: this.y - 40 - 15 - this.bubbleHeight,
        }
      }

      let x

      switch (this.xDirection) {
        case 'left':
          x = this.x - this.bubbleWidth - 20 - 15 - 15
          break
        case 'right':
          x = this.x + this.width + 20 + 15 + 15
          break
      }

      let y = this.y + this.height / 2 - 15 - 15

      return { x, y }
    },

    // eslint-disable-next-line vue/return-in-computed-property
    bubblePath(): string {
      if (this.centerTopAligned) {
        return `
          M${this.x + this.width / 2} ${this.y - 20}
          l-15 -15
          l${-(this.bubbleWidth / 2) + 15} 0
          q-15 0 -15 -15
          l0 -${this.bubbleHeight}
          q0 -15 15 -15
          l${this.bubbleWidth} 0
          q15 0 15 15
          l0 ${this.bubbleHeight}
          q0 15 -15 15
          l${-(this.bubbleWidth / 2) + 15} 0
          l-15 15
        `
      }
      switch (this.xDirection) {
        case 'left':
          return `
            M${this.x - 20} ${this.y + this.height / 2}
            l-15 -15
            l0 -10
            q0 -15 -15 -15
            l-${this.bubbleWidth} 0
            q-15 0 -15 15
            l0 ${this.bubbleHeight}
            q0 15 15 15
            l${this.bubbleWidth} 0
            q15 0 15 -15
            l0 ${-(this.bubbleHeight - 40)}
            l15 -15
          `
        case 'right':
          return `
            M${this.x + this.width + 20} ${this.y + this.height / 2}
            l15 -15
            l0 -10
            q0 -15 15 -15
            l${this.bubbleWidth} 0
            q15 0 15 15
            l0 ${this.bubbleHeight}
            q0 15 -15 15
            l${-this.bubbleWidth} 0
            q-15 0 -15 -15
            l0 ${-(this.bubbleHeight - 40)}
            l-15 -15
          `
      }
    },
  },
  watch: {
    currentStep() {
      this.align()
    },
  },
  mounted() {
    window.addEventListener('resize', this.onResize)
    this.interval = setInterval(this.align, 1000)
  },
  unmounted() {
    window.removeEventListener('resize', this.onResize)
    clearInterval(this.interval)
  },
  methods: {
    onResize() {
      this.screenWidth = window.innerWidth
      this.screenHeight = window.innerHeight
      this.align()
    },
    align() {
      if (this.currentStep && this.currentStep.highlightedElement.children.length > 0) {
        Object.assign(this, calculateChildrenRect(this.currentStep.highlightedElement))
      }
    },
    previous() {
      this.steps.previous()
    },
    next() {
      this.steps.next()
    },
    cancel() {
      this.steps.cancel()
    },
  },
})
</script>

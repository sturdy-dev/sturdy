import type { Step } from './Step'
import { ref } from 'vue'
import type { Ref } from 'vue'

export class StepQueue {
  readonly #steps: Ref<Step[]> = ref([])
  readonly #currentIndex = ref(0)
  readonly #completedStepIds = new Set<string>()
  readonly #onCompleteStep: (stepId: string) => void

  constructor(onCompleteStep: (stepId: string) => void) {
    this.#onCompleteStep = onCompleteStep
  }

  get currentStep(): Step | undefined {
    return this.#steps.value[this.#currentIndex.value]
  }

  get refs(): Record<string, Ref<unknown>> {
    return {
      ['@@stepQueue.steps']: this.#steps,
      ['@@stepQueue.currentIndex']: this.#currentIndex,
    }
  }

  get length() {
    return this.#steps.value.length
  }

  get progress() {
    return this.#currentIndex.value
  }

  get remainingSteps() {
    return this.#steps.value.slice(this.#currentIndex.value)
  }

  get completedSteps() {
    return this.#steps.value.slice(0, this.#currentIndex.value)
  }

  indexOfStep(stepId: string): number | undefined {
    const index = this.#steps.value.findIndex((s) => s?.id === stepId)
    if (index === -1) {
      return undefined
    }
    return index
  }

  sortSteps(steps: Step[]): Step[] {
    const originalSteps = new Map(steps.map((s) => [s.id, s]))
    const stepsWithoutDanglingEdges = steps.map((step) => {
      return {
        id: step.id,
        dependencies: new Set(
          Array.from(step.dependencies ?? []).filter((depId) => steps.some((s) => s.id === depId))
        ),
      }
    })
    const stepsWithDependencies = stepsWithoutDanglingEdges.filter((s) => s.dependencies.size > 0)
    const stepsWithoutDependencies = stepsWithoutDanglingEdges.filter(
      (s) => s.dependencies.size === 0
    )

    const sortedIds: string[] = []
    while (stepsWithoutDependencies.length > 0) {
      const n = stepsWithoutDependencies.shift()!
      sortedIds.push(n.id)
      for (const m of stepsWithDependencies) {
        if (m.dependencies.has(n.id)) {
          m.dependencies.delete(n.id)
          if (m.dependencies.size === 0) {
            stepsWithoutDependencies.push(m)
          }
        }
      }
    }

    return sortedIds.map((l) => originalSteps.get(l)!)
  }

  registerStep(step: Step) {
    if (this.isCompleted(step.id)) {
      return
    }

    this.#steps.value = [...this.completedSteps, ...this.sortSteps([...this.remainingSteps, step])]
  }

  unregisterStep(step: Step) {
    const index = this.indexOfStep(step.id)
    if (index == null) {
      return
    }
    if (index < this.#currentIndex.value) {
      this.previous()
    }
    this.#steps.value = this.#steps.value.slice(0, index).concat(this.#steps.value.slice(index + 1))
  }

  next() {
    if (this.currentStep) {
      this.complete(this.currentStep.id)
    }
    if (this.currentStep && this.isCompleted(this.currentStep.id)) {
      this.#currentIndex.value++
    }
  }

  previous() {
    this.#currentIndex.value--
  }

  cancel() {
    for (const step of this.#steps.value) {
      this.complete(step.id)
    }
    this.#currentIndex.value = this.#steps.value.length
  }

  complete(stepId: string) {
    if (this.isCompleted(stepId)) {
      return
    }

    this.#completedStepIds.add(stepId)
    this.#onCompleteStep(stepId)
  }

  isCompleted(stepId: string): boolean {
    return this.#completedStepIds.has(stepId)
  }
}

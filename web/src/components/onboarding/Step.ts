import type { Slot } from 'vue'

export interface Step {
  id: string
  dependencies?: Set<string>
  highlightedElement: HTMLElement
  title?: Slot
  description?: Slot
}

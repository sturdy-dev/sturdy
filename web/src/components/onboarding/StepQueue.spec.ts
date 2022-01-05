/**
 * @jest-environment jsdom
 */

import { StepQueue } from './StepQueue'
import { Step } from './Step'

function makeStep(id = Math.random().toString(16)): Step {
  const highlightedElement = document.createElement('div')
  highlightedElement.textContent = id
  return {
    id,
    highlightedElement,
    description: () => [],
    title: () => [],
  }
}

function makeQueue() {
  return new StepQueue(() => {
    // noop
  })
}

describe('StepQueue', () => {
  it('starts out empty', () => {
    const queue = makeQueue()
    expect(queue.currentStep).toBeUndefined()
  })

  it('can register a step and select it', () => {
    const queue = makeQueue()

    const step = makeStep()

    queue.registerStep(step)

    expect(queue.currentStep).toEqual(step)

    queue.next()

    expect(queue.currentStep).toBeUndefined()

    queue.previous()

    expect(queue.currentStep).toEqual(step)
  })

  it('can deregister a step after the current one', () => {
    const a = makeStep('a')
    const b = makeStep('b')
    const c = makeStep('c')

    const queue = makeQueue()
    queue.registerStep(a) // <a>
    queue.registerStep(b) // <a>, b
    queue.registerStep(c) // <a>, b, c

    queue.unregisterStep(b) // <a>, c

    expect(queue.currentStep).toEqual(a)
    queue.next() // a, <c>
    expect(queue.currentStep).toEqual(c)
  })

  it('can deregister a step before the current one', () => {
    const a = makeStep('a')
    const b = makeStep('b')
    const c = makeStep('c')

    const queue = makeQueue()
    queue.registerStep(a) // <a>
    queue.registerStep(b) // <a>, b
    queue.registerStep(c) // <a>, b, c

    queue.next() // a, <b>, c

    queue.unregisterStep(a) // <b>, c

    expect(queue.currentStep).toEqual(b)
    queue.next()
    expect(queue.currentStep).toEqual(c)
  })

  it('can deregister the current step', () => {
    const a = makeStep('a')
    const b = makeStep('b')
    const c = makeStep('c')

    const queue = makeQueue()
    queue.registerStep(a) // <a>
    queue.registerStep(b) // <a>, b
    queue.registerStep(c) // <a>, b, c

    queue.next() // a, <b>, c

    queue.unregisterStep(b) // a, <c>

    expect(queue.currentStep).toEqual(c)
  })

  it('can register a step with dependency', () => {
    const a = makeStep('a')
    a.dependencies = new Set(['b'])
    const b = makeStep('b')
    const c = makeStep('c')

    const queue = makeQueue()
    queue.registerStep(a) // <a>
    queue.registerStep(b) // <b>, a
    queue.registerStep(c) // <b>, a, c

    expect(queue.currentStep).toEqual(b)
    queue.next()
    expect(queue.remainingSteps).toContainEqual(a)
    expect(queue.remainingSteps).toContainEqual(c)
  })

  it('can register a step with dependent step that has already passed', () => {
    const a = makeStep('a')
    a.dependencies = new Set(['c'])
    const b = makeStep('b')
    const c = makeStep('c')

    const queue = makeQueue()
    queue.registerStep(a) // <a>
    queue.registerStep(b) // <a>, b
    queue.next() // a, <b>
    queue.registerStep(c) // a, <c>, b

    expect(queue.remainingSteps).toContainEqual(c)
    expect(queue.remainingSteps).toContainEqual(b)
  })

  it('reorders to satisfy dependencies', () => {
    const a = makeStep('a')
    a.dependencies = new Set(['c'])
    const b = makeStep('b')
    const c = makeStep('c')
    c.dependencies = new Set(['b'])

    const queue = makeQueue()
    queue.registerStep(a) // <a>
    queue.registerStep(b) // <a>, b
    queue.registerStep(c) // <b>, c, a

    expect(queue.currentStep).toEqual(b)
    queue.next()
    expect(queue.currentStep).toEqual(c)
    queue.next()
    expect(queue.currentStep).toEqual(a)
  })
})

import type { Meta, Story } from '@storybook/vue3'
import { defineComponent } from 'vue'

import Button from './Button.vue'

export default {
  title: 'Button',
  component: Button,
  argTypes: {
    default: {
      control: 'text',
      description: 'Slot content',
      defaultValue: 'Click me!',
    },
    color: { control: { type: 'select', options: ['white', 'blue'] } },
  },
} as Meta

export const Default: Story = (args) =>
  defineComponent({
    components: { Button },
    setup: () => ({ args }),
    template: '<Button v-bind="args">{{ args.default }}</Button>',
  })

export const LongerText = Default.bind({})
LongerText.args = {
  default: 'Lorem ipsum and so on',
}

export const Blue = Default.bind({})
Blue.args = {
  color: 'blue',
}

import { Meta, Story } from '@storybook/vue3'
import { defineComponent } from 'vue'

import Banner from './Banner.vue'

export default {
  title: 'Banner',
  component: Banner,
  argTypes: {
    message: {
      control: 'text',
      defaultValue: 'Hello!',
    },
    status: { control: { type: 'select', options: ['success', 'info', 'error'] } },
  },
} as Meta

export const Default: Story = (args) =>
  defineComponent({
    components: { Banner },
    setup: () => ({ args }),
    template: '<Banner v-bind="args" />',
  })

export const Success = Default.bind({})
Success.args = {
  status: 'success',
  message: 'Hello world!',
}

export const Info = Default.bind({})
Info.args = {
  status: 'info',
  message: 'Hello world!',
}

export const Error = Default.bind({})
Error.args = {
  status: 'error',
  message: 'Hello world!',
}

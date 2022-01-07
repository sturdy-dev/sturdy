import { createApp } from './main.ts'

createApp(false).then(({ app }) => {
  app.mount('#app')
})

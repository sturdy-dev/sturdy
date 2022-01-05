<template>
  <div
    v-for="(codebase, codebaseIdx) in codebases"
    :key="codebaseIdx"
    class="relative animate-pulse"
  >
    <div
      class="text-gray-900 hover:text-gray-500 flex items-center pl-3 pr-1 py-2 text-sm font-medium rounded-md mx-1 justify-between transition whitespace-nowrap space-x-2 h-10"
    >
      <div class="h-4 bg-gray-300 rounded" :class="codebase.width"></div>

      <button
        class="p-1 hover:bg-warmgray-100 rounded-md cursor-pointer transition flex-shrink-0"
        disabled
      >
        <PlusSmIcon class="w-5 h-5 hover:text-gray-900" />
      </button>
    </div>
    <div class="flex flex-col overflow-hidden text-ellipsis mb-4">
      <template v-for="(workspace, workspaceIdx) in codebase.workspaces" :key="workspaceIdx">
        <div
          class="whitespace-nowrap text-gray-500 text-sm font-medium py-2 pl-2 pr-2 inline-flex items-center relative rounded-md my-0.5 mx-1 group"
        >
          <div
            class="rounded-full h-5 w-5 flex items-center justify-center mr-2 flex-shrink-0 bg-gray-300"
          ></div>
          <div class="h-4 bg-gray-300 rounded" :class="workspace.width"></div>
        </div>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
import { PlusSmIcon } from '@heroicons/vue/solid'
import { PropType } from 'vue'

type placeholder = {
  width: string
}

type workspacePlaceholder = placeholder

type codebasePlaceholder = placeholder & {
  workspaces: workspacePlaceholder[]
}

const widths = ['w-1/4', 'w-2/4', 'w-3/4', 'w-1/2', 'w-full']

const randomWidth = () => widths[Math.floor(Math.random() * widths.length)]

const randomWorkspace = (): workspacePlaceholder => ({
  width: randomWidth(),
})

const randomInt = (min: number, max: number) => {
  min = Math.ceil(min)
  max = Math.floor(max)
  return Math.floor(Math.random() * (max - min + 1)) + min
}

const randomWorkspaces = (): workspacePlaceholder[] => {
  const workspaces = []
  for (let i = 0; i < randomInt(1, 5); i++) {
    workspaces.push(randomWorkspace())
  }
  return workspaces
}

const randomCodebase = (): codebasePlaceholder => ({
  width: randomWidth(),
  workspaces: randomWorkspaces(),
})

const randomCodebases = (): codebasePlaceholder[] => {
  const codebases = []
  for (let i = 0; i < randomInt(2, 5); i++) {
    codebases.push(randomCodebase())
  }
  return codebases
}

export default {
  props: {
    codebases: {
      type: Array as PropType<codebasePlaceholder[]>,
      default: randomCodebases,
    },
  },
  components: {
    PlusSmIcon,
  },
}
</script>

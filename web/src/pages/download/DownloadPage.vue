<template>
  <div>
    <div class="px-6 py-12 bg-gray-800 text-white flex flex-col items-center">
      <h1 class="text-4xl lg:text-6xl font-bold text-center text-yellow-400">Download Sturdy</h1>

      <p class="mt-5 max-w-lg text-center text-gray-200">
        The Sturdy app connects your computer to the cloud, continuously saving your changes. Use
        multiple workspaces to organize your work, and iteratively release your code.
      </p>
    </div>

    <div class="px-6 py-8 flex flex-col items-center">
      <select
        v-model="selectedOsId"
        class="rounded-md bg-gray-50 border border-gray-300 shadow-md py-2 px-3 pr-10"
      >
        <option value="undefined" disabled>Choose your operating system</option>
        <option
          v-for="os in operatingSystems"
          :key="os.id"
          :value="os.id"
          :disabled="os.comingSoon"
        >
          {{ os.name }} {{ os.comingSoon ? '– Coming Soon' : '' }}
        </option>
      </select>
    </div>

    <component
      :is="selectedOs.component"
      v-if="selectedOs?.component"
      class="flex flex-row justify-center items-center gap-2 flex-wrap px-5 mb-7"
    />

    <ul
      v-if="selectedOs != null"
      class="flex flex-row justify-center items-center gap-2 flex-wrap px-5 mb-7"
    >
      <li v-for="dl in selectedOs.archDownloads" :key="dl.id">
        <a
          :href="dl.url"
          type="button"
          download
          class="appearance-none flex gap-1 items-center rounded-md bg-gradient-to-b from-green-500 to-green-600 border border-green-600 text-green-50 font-medium shadow-md py-2 px-3 hover:from-green-400 hover:to-green-500 hover:border-green-500"
        >
          <DownloadIcon class="h-5 w-5" />
          Sturdy for {{ selectedOs.name }} on {{ dl.name }}
        </a>
      </li>
    </ul>

    <img alt="The Sturdy Application" src="./app.png" class="lg:max-w-4xl m-auto w-full" />
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'
import { DownloadIcon } from '@heroicons/vue/solid'
import DownloadLinux from './DownloadLinux.vue'

interface OperatingSystem {
  id: string
  name: string
  archDownloads: ArchDownload[]
  component?: any
}

interface ArchDownload {
  id: string
  name: string
  url: string
}

const operatingSystems: OperatingSystem[] = []

const url = new URL(`https://autoupdate.getsturdy.com/client/`)

// MacOS
const macOs: OperatingSystem = {
  id: 'darwin',
  name: 'MacOS',
  archDownloads: [],
}

// Apple Silicon
macOs.archDownloads.push({
  id: 'darwin-arm64',
  name: 'Apple Silicon',
  url: new URL('darwin/arm64/Install Sturdy.dmg', url).href,
})

// Intel
macOs.archDownloads.push({
  id: 'darwin-amd64',
  name: 'Intel',
  url: new URL('darwin/amd64/Install Sturdy.dmg', url).href,
})

operatingSystems.push(macOs)

// Windows
const windows: OperatingSystem = {
  id: 'windows',
  name: 'Windows',
  archDownloads: [],
}

// x86_64
windows.archDownloads.push({
  id: 'windows-amd64',
  name: 'x86 (64 bit)',
  url: new URL('windows/amd64/Sturdy-Installer.exe', url).href,
})

operatingSystems.push(windows)

// Linux
const linux: OperatingSystem = {
  id: 'linux',
  name: 'Linux',
  archDownloads: [],
  component: DownloadLinux,
}

operatingSystems.push(linux)

export default defineComponent({
  components: {
    DownloadIcon,
  },

  setup() {
    const selectedOsId = ref<undefined | string>()

    if (!import.meta.env.SSR) {
      if (navigator.userAgent.includes('Win')) {
        selectedOsId.value = 'windows'
      } else if (navigator.userAgent.includes('Mac')) {
        selectedOsId.value = 'darwin'
      } else if (navigator.userAgent.includes('Linux')) {
        selectedOsId.value = 'linux'
      }
    }

    return {
      operatingSystems,
      selectedOsId,
    }
  },

  computed: {
    selectedOs() {
      return this.selectedOsId && this.operatingSystems.find((os) => os.id === this.selectedOsId)
    },
  },
})
</script>

<template>
  <div class="divide-y divide-gray-200">
    <div class="pb-5">
      <div class="relative">
        <h3 v-if="downloadLinks[selectedOS]" class="text-sm font-semibold text-gray-800">
          Not using {{ downloadLinks[selectedOS].osName }}?
        </h3>
        <p class="mt-1 text-sm text-gray-600 line-clamp-2">
          Go to instructions for
          <a
            v-if="setupOS !== 'darwin-amd64'"
            href="#"
            class="text-black font-bold"
            @click.stop.prevent="selectedOS = 'darwin-amd64'"
            >macOS</a
          >

          <span v-if="setupOS !== 'linux-amd64'"> or </span>
          <a
            v-if="setupOS !== 'linux-amd64'"
            href="#"
            class="text-black font-bold"
            @click.stop.prevent="selectedOS = 'linux-amd64'"
            >Linux</a
          >

          <span v-if="setupOS !== 'windows-amd64'"> or </span>
          <a
            v-if="setupOS !== 'windows-amd64'"
            href="#"
            class="text-black font-bold"
            @click.stop.prevent="selectedOS = 'windows-amd64'"
            >Windows</a
          >
        </p>
      </div>
    </div>

    <div class="py-5">
      <div class="relative">
        <h3 class="text-sm font-semibold text-gray-800">
          Install Sturdy for {{ downloadLinks[setupOS].osName }}
        </h3>

        <template v-if="downloadLinks[setupOS].showHomebrew">
          <p class="mt-4 text-sm text-gray-600">
            Install the Sturdy client with
            <a class="text-black font-bold" href="https://brew.sh/">Homebrew</a>
            <span v-if="downloadLinks[setupOS].showHomebrewLinuxExtra">, yes, even on Linux!</span>
          </p>
          <div>
            <div class="mt-1 flex rounded-md shadow-sm">
              <div class="relative flex items-stretch flex-grow focus-within:z-10">
                <input
                  id="brew-command"
                  :value="homebrewCommand"
                  type="text"
                  readonly
                  class="focus:ring-blue-500 focus:border-blue-500 block w-full rounded-none rounded-l-md sm:text-sm border-gray-300"
                />
              </div>
              <button
                class="-ml-px relative inline-flex items-center space-x-2 px-4 py-2 border border-gray-300 text-sm font-medium rounded-r-md text-gray-700 bg-gray-50 hover:bg-gray-100 focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
                @click="copyToClipboard('brew-command')"
              >
                <ClipboardCopyIcon class="h-5 w-5 text-gray-400" />
                <span>Copy</span>
              </button>
            </div>
          </div>
          <p class="mt-4 text-sm text-gray-600 line-clamp-2">
            Alternatively, download Sturdy directly and add the binaries (sturdy and sturdy-sync) to
            your
            <span class="font-mono">$PATH</span>
          </p>
        </template>
        <template v-else-if="downloadLinks[setupOS].showWindowsInstructions">
          <p class="mt-4 text-sm text-gray-600">
            Install Sturdy by running this command on the command line
          </p>
          <p class="mt-1 text-sm text-gray-500">Run in PowerShell or the Command Prompt</p>
          <div>
            <div class="mt-1 flex rounded-md shadow-sm">
              <div class="relative flex items-stretch flex-grow focus-within:z-10">
                <input
                  id="powershell-command"
                  :value="powerShellCommand"
                  type="text"
                  readonly
                  class="focus:ring-blue-500 focus:border-blue-500 block w-full rounded-none rounded-l-md sm:text-sm border-gray-300"
                />
              </div>
              <button
                class="-ml-px relative inline-flex items-center space-x-2 px-4 py-2 border border-gray-300 text-sm font-medium rounded-r-md text-gray-700 bg-gray-50 hover:bg-gray-100 focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
                @click="copyToClipboard('powershell-command')"
              >
                <ClipboardCopyIcon class="h-5 w-5 text-gray-400" />
                <span>Copy</span>
              </button>
            </div>
          </div>
          <p class="mt-4 text-sm text-gray-600">
            After installing, restart the Terminal or Command prompt to reload the $PATH.
          </p>
        </template>

        <div v-if="downloadLinks[setupOS].showDirectDownload" class="mt-4 flex space-x-3">
          <a
            :href="downloadLinks[setupOS].location"
            class="inline-flex justify-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-900"
          >
            <DownloadIcon class="-ml-1 mr-2 h-5 w-5 text-gray-800" />
            <span>Direct Download for {{ downloadLinks[setupOS].name }}</span>
          </a>
        </div>

        <div v-if="downloadLinks[setupOS].defaultPath" class="mt-4 font-mono text-xs">
          curl -o sturdy.tar.gz {{ downloadLinks[setupOS].location }}<br />
          tar xzvf sturdy.tar.gz<br />
          mv sturdy{,-sync} {{ downloadLinks[setupOS].defaultPath }}
        </div>
      </div>
    </div>

    <div v-if="downloadLinks[setupOS].showManualInstallationSection" class="py-5">
      <div class="relative">
        <h3 class="text-sm font-semibold text-gray-800">Manual installation instructions</h3>

        <p class="mt-4 text-sm text-gray-600 line-clamp-2">
          If you can't use the installer based installation above, use the following instructions to
          install Sturdy manually.
        </p>

        <div class="mt-4 flex space-x-3">
          <a
            :href="downloadLinks[setupOS].location"
            class="inline-flex justify-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-900"
          >
            <DownloadIcon class="-ml-1 mr-2 h-5 w-5 text-gray-800" />
            <span>Direct Download for {{ downloadLinks[setupOS].name }}</span>
          </a>
        </div>

        <ol class="list-decimal list-inside mt-4 text-sm text-gray-600">
          <li>Download the latest release as a ZIP-file.</li>
          <li>Extract the ZIP archive</li>
          <li>
            Add the directory with the sturdy and sturdy-sync applications to your
            <span class="font-mono">PATH</span> (<a
              href="https://www.architectryan.com/2018/03/17/add-to-the-path-on-windows-10/"
              target="_blank"
              class="text-black inline-flex items-center"
              ><span>instructions</span> <ExternalLinkIcon class="w-4 h-4 ml-1" /> </a
            >)
          </li>
          <li>
            Run Sturdy from
            <a
              class="text-black inline-flex items-center"
              target="_blank"
              href="https://aka.ms/terminal"
              >Windows Terminal
              <ExternalLinkIcon class="w-4 h-4 ml-1" />
            </a>
            (recommended), PowerShell, or the Command Prompt!
          </li>
          <li>
            <i>Optional:</i> Verify your installation by executing
            <span class="font-mono">sturdy version</span>
          </li>
        </ol>
      </div>
    </div>

    <div class="pt-5">
      <div class="relative">
        <h3 class="text-sm font-semibold text-gray-800">Direct downloads</h3>

        <p class="mt-4 text-sm text-gray-600 line-clamp-2">
          If you can't find your OS above (or are using arm64!), grab the download from one of these
          direct links.
        </p>

        <div class="mt-2 flex space-x-3 text-black text-sm font-semibold">
          <ul class="ml-4 list-disc">
            <li v-for="dl in directDownloadLinks" :key="dl.location">
              <a :href="dl.location">{{ dl.name }}</a>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { DownloadIcon, ExternalLinkIcon } from '@heroicons/vue/solid'
import { ClipboardCopyIcon } from '@heroicons/vue/outline'

const latestVersion = 'v0.8.1-beta2'

const downloadLinks = {
  'darwin-amd64': {
    name: 'MacOS (Intel)',
    osName: 'macOS',
    showHomebrew: true,
    showDirectDownload: true,
    defaultPath: '/usr/local/bin',
    location: 'https://getsturdy.com/client/sturdy-' + latestVersion + '-darwin-amd64.tar.gz',
  },
  'darwin-arm64': {
    name: 'MacOS (Apple Silicon)',
    osName: 'macOS',
    showHomebrew: true,
    showDirectDownload: true,
    defaultPath: '/usr/local/bin',
    location: 'https://getsturdy.com/client/sturdy-' + latestVersion + '-darwin-arm64.tar.gz',
  },
  'linux-amd64': {
    name: 'Linux (amd64)',
    osName: 'Linux',
    showHomebrew: true,
    showHomebrewLinuxExtra: true,
    showDirectDownload: true,
    defaultPath: '/usr/bin',
    location: 'https://getsturdy.com/client/sturdy-' + latestVersion + '-linux-amd64.tar.gz',
  },
  'linux-arm64': {
    name: 'Linux (arm64)',
    osName: 'Linux',
    showHomebrew: true,
    showHomebrewLinuxExtra: true,
    showDirectDownload: true,
    defaultPath: '/usr/bin',
    location: 'https://getsturdy.com/client/sturdy-' + latestVersion + '-linux-arm64.tar.gz',
  },
  'windows-amd64': {
    name: 'Windows (amd64)',
    osName: 'Windows',
    showWindowsInstructions: true,
    showDirectDownload: false,
    showManualInstallationSection: true,
    location: 'https://getsturdy.com/client/sturdy-' + latestVersion + '-windows-amd64.zip',
    // installerLocation:
    //   'https://getsturdy.com/client/sturdy-' + latestVersion + '-windows-amd64.msi',
  },
}

const directDownloadLinks = [
  {
    name: 'MacOS (Intel)',
    location: 'https://getsturdy.com/client/sturdy-' + latestVersion + '-darwin-amd64.tar.gz',
  },
  {
    name: 'MacOS (Apple Silicon)',
    location: 'https://getsturdy.com/client/sturdy-' + latestVersion + '-darwin-arm64.tar.gz',
  },
  {
    name: 'Linux (amd64)',
    location: 'https://getsturdy.com/client/sturdy-' + latestVersion + '-linux-amd64.tar.gz',
  },
  {
    name: 'Linux (ARMv8, 64-bit)',
    location: 'https://getsturdy.com/client/sturdy-' + latestVersion + '-linux-arm64.tar.gz',
  },
  {
    name: 'Linux (ARMv7, 32-bit)',
    location: 'https://getsturdy.com/client/sturdy-' + latestVersion + '-linux-arm7.tar.gz',
  },
  {
    name: 'Linux (ARMv6, 32-bit)',
    location: 'https://getsturdy.com/client/sturdy-' + latestVersion + '-linux-arm6.tar.gz',
  },
  {
    name: 'Linux (ARMv5, 32-bit)',
    location: 'https://getsturdy.com/client/sturdy-' + latestVersion + '-linux-arm5.tar.gz',
  },
  // {
  //   name: 'Windows (amd64) - Installer',
  //   location: 'https://getsturdy.com/client/sturdy-' + latestVersion + '-windows-amd64.msi',
  // },
  {
    name: 'Windows (amd64) - Zip',
    location: 'https://getsturdy.com/client/sturdy-' + latestVersion + '-windows-amd64.zip',
  },
]

export default {
  name: 'InstallationInstructions',
  components: { ExternalLinkIcon, ClipboardCopyIcon, DownloadIcon },
  data() {
    return {
      selectedOS: null,
      supportedOS: ['darwin-amd64', 'linux-amd64', 'windows-amd64'],
      downloadLinks: downloadLinks,
      directDownloadLinks: directDownloadLinks,
      powerShellCommand:
        'powershell ". { iwr -useb https://getsturdy.com/client/windows-installer.ps1 } | iex;"',
      homebrewCommand: 'brew install sturdy-dev/tap/sturdy',
    }
  },
  computed: {
    setupOS() {
      return this.selectedOS ? this.selectedOS : this.detectedOS
    },
    detectedOS() {
      let name = 'darwin-amd64'
      if (!import.meta.env.SSR) {
        if (navigator.appVersion.indexOf('Win') !== -1) {
          name = 'windows-amd64'
        } else if (navigator.appVersion.indexOf('Mac') !== -1) {
          name = 'darwin-amd64'
        } else if (navigator.appVersion.indexOf('Linux') !== -1) {
          name = 'linux-amd64'
        }
      }
      return name
    },
  },
  methods: {
    copyToClipboard(el) {
      var copyText = document.getElementById(el)
      copyText.select()
      copyText.setSelectionRange(0, 99999)
      document.execCommand('copy')
    },
  },
}
</script>
<style scoped></style>

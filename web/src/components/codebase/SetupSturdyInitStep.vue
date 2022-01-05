<template>
  <div class="flex rounded-md mb-1 mr-1 w-full bg-green-200">
    <div class="relative flex items-stretch flex-grow focus-within:z-10">
      <input
        id="init-command"
        :value="sturdyInitCommand"
        type="text"
        readonly
        name="codebase_name"
        class="focus:ring-blue-500 focus:border-blue-500 block w-full rounded-none rounded-l-md sm:text-sm border-gray-300"
      />
    </div>
    <button
      class="-ml-px relative inline-flex items-center space-x-2 px-4 py-2 border border-gray-300 text-sm font-medium rounded-r-md text-gray-700 bg-gray-50 hover:bg-gray-100 focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
      @click="copy"
    >
      <svg
        class="h-5 w-5 text-gray-400"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3"
        />
      </svg>
      <span>Copy</span>
    </button>
  </div>
</template>
<script>
export default {
  name: 'SetupSturdyInitStep',
  props: ['codebase', 'userHaveViews'],
  computed: {
    sturdyInitCommand() {
      return 'sturdy init ' + this.codebase?.id + ' ' + this.pathSafeName()
    },
  },
  methods: {
    pathSafeName() {
      if (!this.codebase) {
        return ''
      }
      let safe = this.trim(this.codebase.name.replace(/[^a-z0-9]/gi, '-').toLowerCase(), '-')
      if (safe.length > 0) {
        return safe
      }
      return 'unnamed-codebase'
    },

    copy() {
      let copyText = document.getElementById('init-command')
      copyText.select()
      copyText.setSelectionRange(0, 99999)
      document.execCommand('copy')
    },

    escapeRegex(string) {
      // eslint-disable-next-line
      return string.replace(/[\[\](){}?*+\^$\\.|\-]/g, '\\$&')
    },

    trim(str, characters, flags = 'g') {
      if (typeof str !== 'string' || typeof characters !== 'string' || typeof flags !== 'string') {
        throw new TypeError('argument must be string')
      }

      if (!/^[gi]*$/.test(flags)) {
        throw new TypeError("Invalid flags supplied '" + flags.match(new RegExp('[^gi]*')) + "'")
      }

      characters = this.escapeRegex(characters)

      return str.replace(new RegExp('^[' + characters + ']+|[' + characters + ']+$', flags), '')
    },
  },
}
</script>

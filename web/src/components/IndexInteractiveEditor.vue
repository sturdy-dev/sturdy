<template>
  <div ref="editorTarget" class="relative">
    <div class="mx-auto px-4 sm:px-6 lg:px-8 lg:max-w-7xl flex flex-col items-center">
      <LightIndexPageDifferFile v-if="diffs" :diffs="diffs" class="w-full w-auto" />

      <div class="ml-0 -mt-28 relative md:px-8 lg:px-0 self-center pt-8 lg:pt-0 md:block">
        <div class="mx-auto lg:max-w-2xl xl:max-w-none">
          <div
            class="relative border border-slate-500/80 overflow-hidden rounded-xl shadow-2xl flex bg-light-blue-500 pb-6 md:pb-0"
          >
            <div class="absolute inset-0 bg-slate-800 bg-opacity-80 backdrop-blur-[5px]" />
            <div class="relative w-full flex flex-col">
              <div class="flex-none h-11 flex items-center px-4">
                <div class="flex space-x-1.5">
                  <div class="w-3 h-3 rounded-full bg-red-500" />
                  <div class="w-3 h-3 rounded-full bg-yellow-400" />
                  <div class="w-3 h-3 rounded-full bg-green-400" />
                </div>
              </div>
              <div
                class="relative border-t border-white border-opacity-10 min-h-0 flex-auto flex flex-col"
              >
                <div
                  class="hidden md:block absolute inset-y-0 left-0 bg-slate-900 bg-opacity-25 backdrop-blur-[5px]"
                  style="width: 50px"
                />
                <div class="w-full flex-auto flex min-h-0 overflow-auto">
                  <div class="w-full relative flex-auto">
                    <pre class="flex min-h-full text-xs md:text-sm"><div
                        aria-hidden="true"
                        class="hidden md:block text-white text-opacity-50 text-xs tracking-tight flex-none py-4 pr-4 text-right select-none"
                        style="width:50px"
                    >1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
</div>
                    <code
                        contenteditable
                        class="flex-auto relative block text-slate-200 text-xs tracking-tight pt-4 pb-4 px-4 overflow-auto outline-none"
                        @input="editorupdate"
                    >{{
                        automatedEditorValue
                      }}</code>
                  </pre>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineAsyncComponent, defineComponent, ref } from 'vue'
import Diff from '../diff/diff.js'
import { useElementVisibility } from '@vueuse/core'

export default defineComponent({
  name: 'IndexInteractiveEditor',
  components: {
    LightIndexPageDifferFile: defineAsyncComponent(
      () => import('./differ/LightIndexPageDifferFile.vue')
    ),
  },
  setup() {
    const editorTarget = ref(null)
    const targetIsVisible = useElementVisibility(editorTarget)

    return {
      editorTarget,
      targetIsVisible,
    }
  },
  data: function () {
    return {
      diffs: null,

      seen: false,
      automatedEditorValue: '',
      defaultEditorValue:
        "const http = require('http');\r\n" +
        '\r\n' +
        "const hostname = '127.0.0.1';\r\n" +
        'const port = 3000;\r\n' +
        '\r\n' +
        'const server = http.createServer((req, res) => {\r\n' +
        '  res.statusCode = 200;\r\n' +
        "  res.setHeader('Content-Type', 'text/plain');\r\n" +
        "  res.end('Hello World');\r\n" +
        '});\r\n' +
        '\r\n' +
        'server.listen(port, hostname, () => {\r\n' +
        '  console.log(`http://${hostname}:${port}/`);\r\n' +
        '});\r\n',
    }
  },
  watch: {
    targetIsVisible() {
      if (this.targetIsVisible && !this.seen) {
        this.seen = true
        this.animate()
      }
    },
  },
  mounted() {
    this.automatedEditorValue = this.defaultEditorValue
    this.updateDiffPreview(this.defaultEditorValue)
    // this.animate()
  },
  methods: {
    animate() {
      let rows = this.defaultEditorValue.split(/\r\n|\r|\n/)

      const stringAnimation = (row, prefix, str, suffix) => {
        var res = []
        for (let i = 1; i <= str.length; i++) {
          res.push({
            row: row,
            content: prefix + str.substring(0, i) + suffix,
          })
        }
        return res
      }

      let frames = [].concat(
        stringAnimation(8, "  res.end('Hello ", 'Sturdy!', "');"),
        { row: 5, newRow: true },
        stringAnimation(5, '// ', 'In Sturdy, all changes are streamed live back to', ''),
        { row: 6, newRow: true },
        stringAnimation(6, '// ', 'the cloud. And are ready to review and merge in ', ''),
        { row: 7, newRow: true },
        stringAnimation(7, '// ', 'less than 1 second!', ''),
        { sleep: 3000 },
        { row: 0, newRow: true },
        stringAnimation(0, '// ', 'Psst! This is an interactive editor!', '')
      )

      let fc = 0

      const nextFrame = () => {
        let f = frames[fc]

        if (f.sleep) {
          //
        } else if (f.newRow) {
          if (f.row === 0) {
            rows.unshift('')
          } else {
            rows.splice(f.row, 0, '')
          }
        } else {
          rows[f.row] = f.content
        }

        this.automatedEditorValue = rows.join('\r\n')
        this.updateDiffPreview(this.automatedEditorValue)

        if (fc < frames.length - 1) {
          fc++
          let timeout = f.sleep ? f.sleep : 40
          setTimeout(nextFrame, timeout)
        }
      }

      setTimeout(nextFrame, 200)
    },
    editorupdate(event) {
      this.updateDiffPreview(event.target.innerText)
    },
    updateDiffPreview(value) {
      let d = this.lineDiffFake(this.defaultEditorValue, value + '\r\n')

      // Fallback to default
      if (d === undefined) {
        d = this.defaultEditorValue
      }

      let rows = d.split(/\r\n|\r|\n/).length

      let header =
        'diff --git a/index.js b/index.js\nindex 99e976b..602ccec 100644\n--- a/index.js\n+++ b/index.js\n@@ -1,' +
        rows +
        ' +1,' +
        rows +
        ' @@ // Demo\n'

      this.diffs = {
        orig_name: 'index.js',
        new_name: 'index.js',
        hunks: [{ patch: header + d }],
      }
    },

    appendAllButLast(str, regex, append) {
      const reg = new RegExp(regex, 'g')
      return str.replace(reg, function (match, offset, str) {
        const follow = str.slice(offset)
        const isLast = follow.match(reg).length === 1
        return isLast ? match : match + append
      })
    },

    lineDiffFake(str1, str2) {
      let diff,
        isDiff,
        accumulatedDiff = ''

      let CRS = '(?:\r?\n)'
      diff = Diff.diffLines(str1, str2)

      isDiff = diff.some((item) => {
        return item.added || item.removed
      })

      if (isDiff) {
        diff.forEach((part: any) => {
          let prefix = ''
          if (part.added) {
            prefix = '+'
          } else if (part.removed) {
            prefix = '-'
          } else {
            prefix = ' '
          }
          part.value = this.appendAllButLast(part.value, CRS, prefix)
          part.diff = prefix + part.value
          accumulatedDiff += part.diff
        })
        return accumulatedDiff
      }

      return undefined
    },
  },
})
</script>

<style scoped></style>

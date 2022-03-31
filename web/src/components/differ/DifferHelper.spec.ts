import { getIndicesOf, type SearchableHunk, searchMatches } from './DifferHelper'

describe('DifferHelper', () => {
  it('get indexes length', () => {
    const str1 = `diff --git "a/readme.md" "b/readme.md"
index 48533da..f627f7d 100644
--- "a/readme.md"
+++ "b/readme.md"
@@ -1,1 +1,7 @@
-# testing1
\\ no newline at end of file
+# testing1
+dsada
+dasda
+testing2
+das
\\ no newline at end of file
+testing3
+testing4
\\ no newline at end of file
`
    const str2 =
      '"diff --git "a/package.json" "b/package.json"\n' +
      'index cecf60c..38f28ed 100644\n' +
      '--- "a/package.json"\n' +
      '+++ "b/package.json"\n' +
      '@@ -31,9 +31,9 @@\n' +
      '     "@urql/exchange-graphcache": "~4.3.5",\n' +
      '     "@urql/exchange-retry": "~0.3.0",\n' +
      '     "@urql/introspection": "~0.3.0",\n' +
      '-    "@urql/vue": "~0.6.0",\n' +
      '+    "@urql/vue": "^0.6.0",\n' +
      '     "@vue/apollo-option": "~4.0.0-alpha.11",\n' +
      '-    "@vueuse/core": "~6.4.0",\n' +
      '+    "@vueuse/core": "^6.4.0",\n' +
      '     "@vueuse/head": "~0.7.2",\n' +
      '     "diff2html": "~3.4.13",\n' +
      '     "emoji-js": "~3.6.0",\n' +
      '@@ -44,25 +44,26 @@\n' +
      '     "linkify-string": "~3.0.3",\n' +
      '     "linkifyjs": "~3.0.3",\n' +
      '     "mitt": "~2.1.0",\n' +
      '-    "posthog-js": "~1.17.9",\n' +
      '+    "posthog-js": "^1.17.9",\n' +
      '     "prismjs": "^1.27.0",\n' +
      '-    "subscriptions-transport-ws": "~0.9.19",\n' +
      '+    "subscriptions-transport-ws": "^0.9.19",\n' +
      '     "tributejs": "~5.1.3",\n' +
      '     "twitter-widgets": "~2.0.0",\n' +
      '     "urql": "~2.0.5",\n' +
      '-    "vite-ssr": "~0.15.0",\n' +
      '+    "vite-ssr": "^0.15.0",\n' +
      '     "vue": "~3.2.29",\n' +
      '     "vue-prism-editor": "~2.0.0-alpha.2",\n' +
      '     "vue-router": "~4.0.11"\n' +
      '   },\n' +
      '   "devdependencies": {\n' +
      '-    "@graphql-codegen/cli": "2.3.0",\n' +
      '+    "@graphql-codegen/cli": "^2.3.0",\n' +
      '+    "@graphql-codegen/introspection": "^2.1.1",\n' +
      '     "@graphql-codegen/introspection": "^2.1.1",\n' +
      '     "@graphql-codegen/near-operation-file-preset": "~2.2.0",\n' +
      '-    "@graphql-codegen/typescript": "2.4.0",\n' +
      '-    "@graphql-codegen/typescript-operations": "~2.2.0",\n' +
      '-    "@graphql-codegen/typescript-vue-urql": "~2.2.0",\n' +
      '-    "@graphql-codegen/urql-introspection": "2.1.0",\n' +
      '+    "@graphql-codegen/typescript": "^2.4.0",\n' +
      '+    "@graphql-codegen/typescript-operations": "^2.2.0",\n' +
      '+    "@graphql-codegen/typescript-vue-urql": "^2.2.0",\n' +
      '+    "@graphql-codegen/urql-introspection": "^2.1.0",\n' +
      '     "@headlessui/vue": "~1.4.0",\n' +
      '     "@heroicons/vue": "~1.0.3",\n' +
      '     "@tailwindcss/forms": "~0.4.0",\n' +
      '@@ -71,41 +72,41 @@\n' +
      '     "@testing-library/jest-dom": "~5.15.1",\n' +
      '     "@types/emoji-js": "~3.5.0",\n' +
      '     "@types/jest": "~27.0.3",\n' +
      '-    "@typescript-eslint/eslint-plugin": "~5.8.1",\n' +
      '-    "@typescript-eslint/parser": "~5.8.1",\n' +
      '+    "@typescript-eslint/eslint-plugin": "^5.8.1",\n' +
      '+    "@typescript-eslint/parser": "^5.8.1",\n' +
      '     "@vitejs/plugin-vue": "~2.1.0",\n' +
      '-    "@vitejs/plugin-vue-jsx": "~1.3.0",\n' +
      '     "@vue/cli-plugin-babel": "~4.5.15",\n' +
      '     "@vue/compiler-sfc": "~3.2.23",\n' +
      '-    "@vue/eslint-config-typescript": "~9.1.0",\n' +
      '+    "@vue/eslint-config-typescript": "^9.1.0",\n' +
      '     "@vue/test-utils": "~2.0.0-rc.17",\n' +
      '+    "@vue/tsconfig": "^0.1.3",\n' +
      '     "@vue/vue3-jest": "~27.0.0-alpha.3",\n' +
      '-    "@vuedx/typecheck": "~0.7.4",\n' +
      '-    "@vuedx/typescript-plugin-vue": "~0.7.4",\n' +
      '+    "@vuedx/typecheck": "^0.7.4",\n' +
      '+    "@vuedx/typescript-plugin-vue": "^0.7.4",\n' +
      '     "autoprefixer": "~10.4.0",\n' +
      '     "babel-loader": "~8.2.3",\n' +
      '     "core-js": "~3.6.5",\n' +
      '-    "eslint": "~8.5.0",\n' +
      '-    "eslint-config-prettier": "~8.3.0",\n' +
      '-    "eslint-plugin-file-progress": "~1.1.1",\n' +
      '-    "eslint-plugin-jest": "~25.3.0",\n' +
      '-    "eslint-plugin-prettier": "~4.0.0",\n' +
      '-    "eslint-plugin-prettier-vue": "~3.1.0",\n' +
      '-    "eslint-plugin-vue": "~8.2.0",\n' +
      '-    "jest": "~27.3.1",\n' +
      '+    "eslint": "^8.5.0",\n' +
      '+    "eslint-config-prettier": "^8.3.0",\n' +
      '+    "eslint-plugin-file-progress": "^1.1.1",\n' +
      '+    "eslint-plugin-jest": "^25.3.0",\n' +
      '+    "eslint-plugin-prettier": "^4.0.0",\n' +
      '+    "eslint-plugin-prettier-vue": "^3.1.0",\n' +
      '+    "eslint-plugin-vue": "^8.2.0",\n' +
      '+    "jest": "^27.3.1",\n' +
      '     "node-fetch": "~2.6.1",\n' +
      '-    "node-sass": "~6.0.0",\n' +
      '+    "node-sass": "^7.0.1",\n' +
      '     "postcss": "~8.4.4",\n' +
      '     "prettier": "~2.5.1",\n' +
      '     "sass": "~1.32.8",\n' +
      '-    "tailwindcss": "~3.0.21",\n' +
      '+    "tailwindcss": "^3.0.21",\n' +
      '     "ts-jest": "~27.0.7",\n' +
      '-    "typescript": "~4.3.0",\n' +
      '+    "typescript": "^4.6.2",\n' +
      '     "vite": "~2.6.14",\n' +
      '     "vite-imagetools": "^4.0.3",\n' +
      '     "vite-plugin-eslint": "~1.3.0",\n' +
      '-    "vue-loader": "~17.0.0",\n' +
      '-    "vue-tsc": "~0.29.6",\n' +
      '+    "vue-loader": "^17.0.0",\n' +
      '+    "vue-tsc": "^0.29.6",\n' +
      '     "yarn-deduplicate": "^3.1.0"\n' +
      '   }\n' +
      ' }\n' +
      '"'

    expect(getIndicesOf('testing', str1, false).length).toEqual(5)
    expect(getIndicesOf('da', str1, false).length).toEqual(3)
    expect(getIndicesOf('sa', str1, false).length).toEqual(1)
    expect(getIndicesOf('f', str1, false).length).toEqual(0)
    expect(getIndicesOf('urql', str2, false).length).toEqual(10)
    expect(getIndicesOf('vue', str2, false).length).toEqual(34)

    {
      const hunk1: SearchableHunk = {
        hunkID: 'd320fd593ace289810cb0991e437fec054ef34d4215bf8e07a429c17fa1fea36',
        patch: str1,
      }
      const map1 = new Map<string, number[]>()
      map1.set(hunk1.hunkID, getIndicesOf('testing', hunk1.patch, false))
      expect(searchMatches(map1, [hunk1]).size).toEqual(5)
    }

    {
      const hunk2: SearchableHunk = {
        hunkID: 'd320fd593ace289810cb0991e437fec054ef34d4215bf8e07a429c17fa1fea36',
        patch: str2,
      }

      const map2 = new Map<string, number[]>()
      map2.set(hunk2.hunkID, getIndicesOf('urql', hunk2.patch, false))
      const matches = searchMatches(map2, [hunk2])
      const rowsIndexes: [string, string][] = []
      matches.forEach((x) => {
        const match = x.split('-')
        rowsIndexes.push([match[2], match[1]])
      })
      expect(rowsIndexes.length).toEqual(10)
      expect(rowsIndexes[0]).toEqual(['0', '0'])
      expect(rowsIndexes[1]).toEqual(['0', '1'])
      expect(rowsIndexes[2]).toEqual(['0', '2'])
      expect(rowsIndexes[3]).toEqual(['0', '3'])
      expect(rowsIndexes[4]).toEqual(['0', '4'])
      expect(rowsIndexes[5]).toEqual(['1', '10'])
      expect(rowsIndexes[6]).toEqual(['1', '25'])
      expect(rowsIndexes[7]).toEqual(['1', '26'])
      expect(rowsIndexes[8]).toEqual(['1', '29'])
      expect(rowsIndexes[9]).toEqual(['1', '30'])
    }
  })
})

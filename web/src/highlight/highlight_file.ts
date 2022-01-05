import hljs from 'highlight.js/lib/core'

import javascript from 'highlight.js/lib/languages/javascript'
import 'highlight.js/styles/github.css'

hljs.registerLanguage('javascript', javascript)

export default (input: string, lang: string): any => {
  if (lang === 'vue') {
    lang = 'jsx'
  }

  if (!hljs.getLanguage(lang)) {
    return false
  }

  return hljs.highlight(input, { language: lang }).value
}

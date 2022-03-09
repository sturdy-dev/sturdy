import hljs from 'highlight.js/lib/core'

import javascript from 'highlight.js/lib/languages/javascript'
import 'highlight.js/styles/github.css'
import type { Block } from '../components/differ/event'
hljs.registerLanguage('javascript', javascript)

export default (input: Block[], language: string, highlightEnabled: boolean) => {
  if (language === 'vue') {
    language = 'jsx'
  }
  const shouldHighlight = highlightEnabled && hljs.getLanguage(language)
  return input.map(({ header, lines }) => ({
    header,
    lines: lines.map(({ oldNumber, newNumber, type, content }) => ({
      oldNumber,
      newNumber,
      type,
      prefix: content.substring(0, 1),
      originalContent: content.substring(1),
      content: shouldHighlight ? hljs.highlight(content.substring(1), { language }).value : '',
    })),
  }))
}

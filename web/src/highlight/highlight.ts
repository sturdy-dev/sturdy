import hljs from 'highlight.js/lib/core'

import javascript from 'highlight.js/lib/languages/javascript'
import 'highlight.js/styles/github.css'
import { Block, HighlightedBlock, HighlightedLine, Line } from '../components/differ/event'
hljs.registerLanguage('javascript', javascript)

export default (
  input: Array<Block>,
  lang: string,
  highlightEnabled: boolean
): Array<HighlightedBlock> => {
  if (lang === 'vue') {
    lang = 'jsx'
  }
  const out = new Array<HighlightedBlock>()
  input.forEach((block: Block) => {
    const hlBlock: HighlightedBlock = {
      header: block.header,
      lines: [],
    }
    block.lines.forEach((line: Line) => {
      const hlLine: HighlightedLine = {
        oldNumber: line.oldNumber,
        newNumber: line.newNumber,
        type: line.type,
        prefix: line.content.substring(0, 1),
        originalContent: line.content.substring(1),
        content: '',
      }
      if (hljs.getLanguage(lang) && highlightEnabled) {
        hlLine.content = hljs.highlight(line.content.substring(1), {
          language: lang,
        }).value
      }
      hlBlock.lines.push(hlLine)
    })
    out.push(hlBlock)
  })
  return out
}

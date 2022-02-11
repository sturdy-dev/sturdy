import { Renderer } from '@ts-stack/markdown'

export class SturdyMarkdownRenderer extends Renderer {
  html(html: string): string {
    const res = super.html(html)
    const parser = new DOMParser()
    const document = parser.parseFromString(res, 'text/html')

    // Convert height and width attributes to style attributes
    // This is necessary to render images as expected, with Tailwind
    // See https://github.com/tailwindlabs/tailwindcss/issues/506 for more background
    const imgs = document.querySelectorAll<HTMLElement>('img[height], img[width]')

    imgs.forEach((img) => {
      const height = img.getAttribute('height')
      if (height) {
        img.style.height = `${height}px`
      }

      const width = img.getAttribute('width')
      if (width) {
        img.style.width = `${width}px`
      }
    })

    return document.documentElement.innerHTML
  }
}

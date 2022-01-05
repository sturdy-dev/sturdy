/* eslint-disable @typescript-eslint/no-var-requires */

// Pre-rendering of static pages
// Pre-rendered pages are served as static HTML (at first), and are later "Hydrated" and re-rendered clientside
const fs = require('fs')
const path = require('path')
const toAbsolute = (p) => path.resolve(__dirname, p)
const manifest = require('./dist/static/ssr-manifest.json')
const template = fs.readFileSync(toAbsolute('dist/static/index.html'), 'utf-8')
const { render, routerRoutes } = require('./dist/server/entry-server.js')
const { renderHeadToString } = require('@vueuse/head')

let routesToPrerender = routerRoutes
  .filter((r) => r.meta?.nonApp && !r.meta?.isAuth && !r.meta?.skipPrerender)
  .map((r) => r.path)

function ensureDirectoryExistence(filePath) {
  const dirname = path.dirname(filePath)
  if (fs.existsSync(dirname)) {
    return true
  }
  ensureDirectoryExistence(dirname)
  fs.mkdirSync(dirname)
}

function buildSitemap(routes) {
  // TODO: Support per-page lastmod
  let lastmod = new Date().toISOString().slice(0, 10)
  let urls = routes
    .map(
      (route) => `<url><loc>https://getsturdy.com${route}</loc><lastmod>${lastmod}</lastmod></url>`
    )
    .join('\n')

  return `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">${urls}</urlset>`
}

;(async () => {
  // pre-render each route...
  for (const url of routesToPrerender) {
    console.log('Rendering:', url)

    const [appHtml, preloadLinks, head] = await render(url, manifest)

    // `head` is created from `createHead()`
    const { headTags, htmlAttrs, bodyAttrs } = renderHeadToString(head)

    let finalHTML = template
      .replace('<div id="app">', '<div id="app" data-server-rendered="true">')
      .replace(`<!--head-tags-->`, headTags)
      .replace(`<!--preload-links-->`, preloadLinks)
      .replace(`<!--app-html-->`, appHtml)
      .replace(`<!--htmlAttrs-->`, htmlAttrs)
      .replace(`<!--bodyAttrs-->`, bodyAttrs)

    const filePath = `dist/static${url === '/' ? '/index' : url}.html`

    ensureDirectoryExistence(filePath)

    fs.writeFileSync(toAbsolute(filePath), finalHTML)
  }

  // Create sitemap
  fs.writeFileSync(toAbsolute('dist/static/sitemap.xml'), buildSitemap(routesToPrerender))

  // done, delete ssr manifest
  fs.unlinkSync(toAbsolute('dist/static/ssr-manifest.json'))

  console.log('~~~ Take these routes and add them to the Lambda@Edge function ~~~')
  console.log(routesToPrerender)
})()

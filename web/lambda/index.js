'use strict'

// This file is tracked in web/lambda/index.js
// Do not edit it on the web editor!

exports.handler = (event, context, callback) => {
  const prerendered = [
    '/',
    '/download',
    '/about',
    '/contact',
    '/privacy',
    '/terms-of-service',
    '/pricing',
    '/security',
    '/docs',
    '/kth',
    '/aoc',
    '/features/access-control',
    '/features/instant-integration',
    '/quickstart',
    '/syncing',
    '/features/instant-switching',
    '/features/integrations',
    '/features/live',
    '/features/workflow',
    '/features/conflicts',
    '/features/large-files',
    '/docs/cli',
    '/api',
    '/docs/access-control',
    '/docs/suggestions',
    '/docs/sturdy-for-git-users',
    '/features/migrate-from-github',
    '/docs/continuous-integration',
    '/blog/2021-12-17-graphql-componentized-uis',
    '/blog/2021-12-07-launching-the-sturdy-app',
    '/blog/2021-11-29-scaling-teams',
    '/blog/2021-11-22-sturdy-the-app-is-coming',
    '/blog/2021-09-29-acls-and-a-fresh-hot-look',
    '/blog/2021-09-09-large-files',
    '/blog/2021-08-18-unbreaking-code-collaboration',
    '/blog/2021-08-12-signup-is-open',
    '/blog/2021-06-10-humane-code-review',
    '/blog/2021-05-06-importing-from-git',
    '/blog/2021-04-16-share-now',
    '/blog/2021-04-01-restore-to-any-point-in-time',
    '/blog/2021-03-24-yc-w21-demo-day',
    '/blog/2021-03-18-this-week-at-sturdy',
    '/blog',
    '/careers',
    '/careers/founding-backend-engineer',
    '/careers/founding-frontend-engineer',
    '/v2/docs',
    '/v2/docs/how-sturdy-augments-git',
    '/v2/docs/working-in-the-open'
  ]

  const request = event.Records[0].cf.request
  // Don't modify request
  if (
    request.uri.startsWith('/css/') ||
    request.uri.startsWith('/img/') ||
    request.uri.startsWith('/js/') ||
    request.uri.startsWith('/client/') ||
    request.uri.startsWith('/assets/') ||
    request.uri.startsWith('/favicon.ico') ||
    request.uri.startsWith('/sitemap.xml') ||
    request.uri.startsWith('/robots.txt')
  ) {
    callback(null, request)
  } else if (request.uri === '/' || request.uri === '/index.html') {
    request.uri = '/index.html'
    callback(null, request)
  } else if (prerendered.includes(request.uri)) {
    request.uri = request.uri + '.html'
    callback(null, request)
  } else {
    request.uri = '/client-side-render.html'
    callback(null, request)
  }
}

'use strict'

// This file is tracked in web/lambda/index.js
// Do not edit it on the web editor!

exports.handler = (event, context, callback) => {
  const prerendered = [
    '/download',
    '/about',
    '/contact',
    '/pricing',
    '/security',
    '/kth',
    '/aoc',
    '/features/access-control',
    '/features/instant-integration',
    '/features/instant-switching',
    '/features/integrations',
    '/features/live',
    '/features/conflicts',
    '/features/large-files',
    '/docs/cli',
    '/api',
    '/docs/access-control',
    '/docs/suggestions',
    '/features/migrate-from-github',
    '/docs/continuous-integration',
    '/blog/2022-04-12-this-week-at-sturdy',
    '/blog/2022-03-10-introducing-draft-changes',
    '/blog/2022-02-21-sturdy-is-now-open-source',
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
    '/careers/senior-backend-software-engineer',
    '/careers/senior-frontend-software-engineer',
    '/careers/full-stack-engineer',
    '/',
    '/docs',
    '/docs/product-intro',
    '/docs/how-sturdy-interacts-with-git',
    '/docs/working-in-the-open',
    '/docs/how-to-ship-software-to-production',
    '/docs/how-to-collaborate-with-others',
    '/docs/how-to-edit-code',
    '/docs/how-to-setup-sturdy-with-github',
    '/docs/how-to-switch-between-tasks',
    '/docs/quickstart',
    '/docs/using-sturdy',
    '/docs/pricing',
    '/docs/index',
    '/docs/self-hosted',
    '/handbook/code-of-conduct',
    '/handbook/releases',
    '/docs/cli/install',
    '/docs/integrations/git',
    '/docs/integrations/git/azure-devops',
    '/docs/integrations/git/gitlab',
    '/privacy',
    '/terms-of-service',
  ]

  const request = event.Records[0].cf.request

  // Redirect go-imports of "getsturdy.com/api" to "github.com/sturdy-dev/sturdy"
  if (
    (request.uri.startsWith('/api') || request.uri.startsWith('/ssh') || request.uri === '/') &&
    request.querystring === 'go-get=1'
  ) {
    const content = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="go-import" content="getsturdy.com git https://github.com/sturdy-dev/sturdy">
    <meta name="go-source" content="getsturdy.com https://github.com/sturdy-dev/sturdy https://github.com/sturdy-dev/sturdy/tree/master{/dir} https://github.com/sturdy-dev/sturdy/blob/master{/dir}/{file}#L{line}">
    <title>Hello Gophers!</title>
  </head>
  <body>
    <p><a href="https://github.com/sturdy-dev/sturdy">https://github.com/sturdy-dev/sturdy</a></p>
  </body>
</html>
`

    const response = {
      status: '200',
      statusDescription: 'OK',
      headers: {
        'cache-control': [
          {
            key: 'Cache-Control',
            value: 'max-age=100',
          },
        ],
        'content-type': [
          {
            key: 'Content-Type',
            value: 'text/html',
          },
        ],
      },
      body: content,
    }
    callback(null, response)
    return
  }

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

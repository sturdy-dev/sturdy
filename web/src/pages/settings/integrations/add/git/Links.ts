export const removeBasicAuth = function (url: string): string {
  if (!url.includes('@')) {
    return url
  }
  const isHTTP = url.startsWith('http://') || url.startsWith('https://')
  if (!isHTTP) {
    return url
  }

  const slashslash = url.indexOf('://')
  const at = url.indexOf('@')

  return url.substring(0, slashslash + 3) + url.substring(at + 1)
}

export const defaultLinkRepo = function (gitURL: string): string | null {
  gitURL = removeBasicAuth(gitURL)

  // GitLab HTTP
  if (gitURL.startsWith('https://gitlab.com/') && gitURL.endsWith('.git')) {
    return gitURL.substring(0, gitURL.length - 4)
  }

  // GitLab SSH
  if (gitURL.startsWith('git@gitlab.com:') && gitURL.endsWith('.git')) {
    return 'https://gitlab.com/' + gitURL.substring(15, gitURL.length - 4)
  }

  // BitBucket HTTP
  if (gitURL.startsWith('https://bitbucket.org/') && gitURL.endsWith('.git')) {
    return gitURL.substring(0, gitURL.length - 4)
  }

  // BitBucket SSH
  if (gitURL.startsWith('git@bitbucket.org:') && gitURL.endsWith('.git')) {
    return 'https://bitbucket.org/' + gitURL.substring(18, gitURL.length - 4)
  }

  // Azure HTTP
  if (gitURL.includes('dev.azure.com') && gitURL.includes('/_git/')) {
    const s = gitURL.indexOf('dev.azure.com/')
    const git = gitURL.indexOf('/_git/')
    const parts = gitURL.substring(s, git).split('/')
    if (parts.length === 3) {
      return 'https://' + parts[0] + '/' + parts[1] + '/_git/' + parts[2]
    }
  }

  // Azure SSH
  if (gitURL.startsWith('git@ssh.dev.azure.com:v3/')) {
    const parts = gitURL.substring(25).split('/')
    return 'https://dev.azure.com/' + parts[0] + '/' + parts[1] + '/_git/' + parts[2]
  }

  return null
}

export const defaultLinkBranch = function (gitURL: string): string | null {
  gitURL = removeBasicAuth(gitURL)

  // GitLab HTTP
  if (gitURL.startsWith('https://gitlab.com/') && gitURL.endsWith('.git')) {
    return gitURL.substring(0, gitURL.length - 4) + '/-/tree/${BRANCH_NAME}'
  }

  // GitLab SSH
  if (gitURL.startsWith('git@gitlab.com:') && gitURL.endsWith('.git')) {
    return (
      'https://gitlab.com/' + gitURL.substring(15, gitURL.length - 4) + '/-/tree/${BRANCH_NAME}'
    )
  }

  // BitBucket HTTP
  if (gitURL.startsWith('https://bitbucket.org/') && gitURL.endsWith('.git')) {
    return gitURL.substring(0, gitURL.length - 4) + '/branch/${BRANCH_NAME}'
  }

  // BitBucket SSH
  if (gitURL.startsWith('git@bitbucket.org:') && gitURL.endsWith('.git')) {
    return (
      'https://bitbucket.org/' + gitURL.substring(18, gitURL.length - 4) + '/branch/${BRANCH_NAME}'
    )
  }

  // Azure HTTP
  if (gitURL.includes('dev.azure.com') && gitURL.includes('/_git/')) {
    const s = gitURL.indexOf('dev.azure.com/')
    const git = gitURL.indexOf('/_git/')
    return (
      'https://' +
      gitURL.substring(s, git) +
      '/_git/' +
      gitURL.substring(git + 6) +
      '?version=GB${BRANCH_NAME}'
    )
  }

  // Azure SSH
  if (gitURL.startsWith('git@ssh.dev.azure.com:v3/')) {
    const parts = gitURL.substring(25).split('/')
    return (
      'https://dev.azure.com/' +
      parts[0] +
      '/' +
      parts[1] +
      '/_git/' +
      parts[2] +
      '?version=GB${BRANCH_NAME}'
    )
  }

  return null
}

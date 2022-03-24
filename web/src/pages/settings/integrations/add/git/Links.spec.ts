import { defaultLinkBranch, defaultLinkRepo, removeBasicAuth } from './Links'

describe('Links', () => {
  it('Remove Basic Auth', () => {
    expect(removeBasicAuth('https://foobar@gitlab.com/zegl/sturdy-push.git')).toEqual(
      'https://gitlab.com/zegl/sturdy-push.git'
    )
    expect(removeBasicAuth('http://foobar@gitlab.com/zegl/sturdy-push.git')).toEqual(
      'http://gitlab.com/zegl/sturdy-push.git'
    )
    expect(removeBasicAuth('http://example.com')).toEqual('http://example.com')
  })

  it('GitLab HTTP', () => {
    const link = 'https://gitlab.com/zegl/sturdy-push.git'
    expect(defaultLinkRepo(link)).toEqual('https://gitlab.com/zegl/sturdy-push')
    expect(defaultLinkBranch(link)).toEqual(
      'https://gitlab.com/zegl/sturdy-push/-/tree/${BRANCH_NAME}'
    )
  })

  it('GitLab SSH', () => {
    const link = 'git@gitlab.com:zegl/sturdy-push.git'
    expect(defaultLinkRepo(link)).toEqual('https://gitlab.com/zegl/sturdy-push')
    expect(defaultLinkBranch(link)).toEqual(
      'https://gitlab.com/zegl/sturdy-push/-/tree/${BRANCH_NAME}'
    )
  })

  it('Azure HTTP', () => {
    const link =
      'https://gustav0446@dev.azure.com/gustav0446/yolo_sturdy_test/_git/yolo_sturdy_test'
    expect(defaultLinkRepo(link)).toEqual('https://dev.azure.com/gustav0446/_git/yolo_sturdy_test')
    expect(defaultLinkBranch(link)).toEqual(
      'https://dev.azure.com/gustav0446/yolo_sturdy_test/_git/yolo_sturdy_test?version=GB${BRANCH_NAME}'
    )
    expect(
      defaultLinkBranch('https://getsturdy@dev.azure.com/getsturdy/sturdy/_git/sturdy-on-azure')
    ).toEqual(
      'https://dev.azure.com/getsturdy/sturdy/_git/sturdy-on-azure?version=GB${BRANCH_NAME}'
    )
  })

  it('Azure SSH', () => {
    const link = 'git@ssh.dev.azure.com:v3/getsturdy/gustav-sturdy-haxx/sturdy-on-azure'
    expect(defaultLinkRepo(link)).toEqual(
      'https://dev.azure.com/getsturdy/gustav-sturdy-haxx/_git/sturdy-on-azure'
    )
    expect(defaultLinkBranch(link)).toEqual(
      'https://dev.azure.com/getsturdy/gustav-sturdy-haxx/_git/sturdy-on-azure?version=GB${BRANCH_NAME}'
    )
    expect(
      defaultLinkBranch('https://getsturdy@dev.azure.com/getsturdy/sturdy/_git/sturdy-on-azure')
    ).toEqual(
      'https://dev.azure.com/getsturdy/sturdy/_git/sturdy-on-azure?version=GB${BRANCH_NAME}'
    )
  })

  it('BitBucket HTTP', () => {
    const link = 'https://zegl@bitbucket.org/zegl/taget-go.git'
    expect(defaultLinkRepo(link)).toEqual('https://bitbucket.org/zegl/taget-go')
    expect(defaultLinkBranch(link)).toEqual(
      'https://bitbucket.org/zegl/taget-go/branch/${BRANCH_NAME}'
    )
  })

  it('BitBucket SSH', () => {
    const link = 'git@bitbucket.org:zegl/taget-go.git'
    expect(defaultLinkRepo(link)).toEqual('https://bitbucket.org/zegl/taget-go')
    expect(defaultLinkBranch(link)).toEqual(
      'https://bitbucket.org/zegl/taget-go/branch/${BRANCH_NAME}'
    )
  })
})

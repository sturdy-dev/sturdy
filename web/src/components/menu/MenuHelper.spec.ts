import { NavigationWorkspace, WorkspaceIndex } from './MenuHelper'

function wsWithName(name: string): NavigationWorkspace {
  return {
    name: name,
    id: name,
    isCurrent: false,
    author: { id: 'yes', avatarUrl: 'foo', name: 'foo' },
    isUnread: false,
    badgeCount: 0,
    presence: [],
    isAuthorCoding: false,
    currentView: undefined,
    isOwnedByUser: false,
    suggestingWorkspaces: [],
  }
}

describe('MenuHelper', () => {
  it('simple workspace index', () => {
    const workspaces = [wsWithName('foo'), wsWithName('bar'), wsWithName('baz')]

    expect(WorkspaceIndex('foo', workspaces)).toEqual(0)
    expect(WorkspaceIndex('bar', workspaces)).toEqual(1)
    expect(WorkspaceIndex('baz', workspaces)).toEqual(2)
  })

  it('with suggestions', () => {
    const bar = wsWithName('bar')

    bar.suggestingWorkspaces = [
      {
        id: 's1-ws',
        isCurrent: false,
        createdAt: 0,
      },
      {
        id: 's2-ws',
        isCurrent: false,
        createdAt: 0,
      },
    ]

    const workspaces = [wsWithName('foo'), bar, wsWithName('baz')]

    expect(WorkspaceIndex('foo', workspaces)).toEqual(0)
    expect(WorkspaceIndex('bar', workspaces)).toEqual(1)
    expect(WorkspaceIndex('s1-ws', workspaces)).toEqual(2)
    expect(WorkspaceIndex('s2-ws', workspaces)).toEqual(3)
    expect(WorkspaceIndex('baz', workspaces)).toEqual(4)
  })

  it('real world data', () => {
    const workspaces = [
      {
        id: 'bb1e2ed4-1369-4b1b-89fa-26faf541c312',
        name: "Fo ofooo's Workspace",
        isCurrent: false,
        author: {
          id: '7a3534d1-13de-4664-8eac-106f4a49664e',
          name: 'Fo ofooo',
          avatarUrl: null,
        },
        isUnread: false,
        badgeCount: 0,
        presence: [],
        isAuthorCoding: true,
        currentView: undefined,
        isOwnedByUser: true,
        suggestingWorkspaces: [],
      },
      {
        id: 'b9006ec0-9e85-4f20-b30f-550f4d47131f',
        name: "Fo ofooo's Workspace",
        isCurrent: true,
        author: {
          id: '7a3534d1-13de-4664-8eac-106f4a49664e',
          name: 'Fo ofooo',
          avatarUrl: null,
        },
        isUnread: false,
        badgeCount: 0,
        presence: [],
        isAuthorCoding: true,

        currentView: undefined,
        isOwnedByUser: true,
        suggestingWorkspaces: [],
      },
      {
        id: '17116e47-8310-4af9-a70f-3cd2fa8ea333',
        name: "Fo ofooo's Workspace",
        isCurrent: false,
        author: {
          id: '7a3534d1-13de-4664-8eac-106f4a49664e',
          name: 'Fo ofooo',
          avatarUrl: null,
        },
        isUnread: false,
        badgeCount: 0,
        presence: [],
        isAuthorCoding: false,

        currentView: undefined,
        isOwnedByUser: true,
        suggestingWorkspaces: [
          {
            id: 's1-ws',
            isCurrent: false,
            createdAt: 0,
          },
          {
            id: 's1-ws',
            isCurrent: false,
            createdAt: 0,
          },
        ],
      },
      {
        id: '5b79fb33-20fd-494a-a7af-87f8af5995e4',
        name: 'Add README',
        isCurrent: false,
        author: {
          id: '7a3534d1-13de-4664-8eac-106f4a49664e',
          name: 'Fo ofooo',
          avatarUrl: null,
        },
        isUnread: false,
        badgeCount: 0,
        presence: [],
        isAuthorCoding: false,
        currentView: undefined,
        isOwnedByUser: true,
        suggestingWorkspaces: [],
      },
      {
        id: 'fad0199b-ac3c-4846-beb9-d686633cb3ce',
        name: 'Fork of Add README',
        isCurrent: false,
        author: {
          id: '7a3534d1-13de-4664-8eac-106f4a49664e',
          name: 'Fo ofooo',
          avatarUrl: null,
        },
        isUnread: false,
        badgeCount: 0,
        presence: [],
        isAuthorCoding: false,
        isOwnedByUser: true,
        suggestingWorkspaces: [],
      },
      {
        id: '16e7fe6d-7187-46bd-bfc1-a63c0f8611bd',
        name: 'Fork of Add README',
        isCurrent: false,
        author: {
          id: '7a3534d1-13de-4664-8eac-106f4a49664e',
          name: 'Fo ofooo',
          avatarUrl: null,
        },
        isUnread: false,
        badgeCount: 0,
        presence: [],
        isAuthorCoding: false,
        isOwnedByUser: true,
        suggestingWorkspaces: [],
      },
    ]

    expect(WorkspaceIndex('17116e47-8310-4af9-a70f-3cd2fa8ea333', workspaces)).toEqual(2) // workspace with suggestions
    expect(WorkspaceIndex('5b79fb33-20fd-494a-a7af-87f8af5995e4', workspaces)).toEqual(5) // workspace after 2 suggestions
    expect(WorkspaceIndex('fad0199b-ac3c-4846-beb9-d686633cb3ce', workspaces)).toEqual(6) // next
  })
})

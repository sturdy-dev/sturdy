import {
  PresenceFragment as WorkspacePresence,
  StackedMenu_ViewFragment as View,
  StackedMenu_WorkspaceAuthorFragment as Author,
} from './__generated__/StackedMenu'

export type NavigationSuggestingWorkspace = {
  id: string
  isCurrent: boolean
  createdAt: number
}

export type NavigationWorkspace = {
  name: string
  id: string
  isCurrent: boolean
  author: Author
  isUnread: boolean
  badgeCount: number
  presence: WorkspacePresence[]
  isAuthorCoding: boolean
  currentView?: View
  isOwnedByUser: boolean
  suggestingWorkspaces: NavigationSuggestingWorkspace[]
  renderNameItalics: boolean
}

export type NavigationView = {
  id: string
  isConnectedToWorkspace: boolean
  workspaceId?: string
  workspaceIndex?: number
  data: View
}

export type NavigationCodebase = {
  id: string
  name: string
  slug: string
  isCurrent: boolean
  workspaces: NavigationWorkspace[]
  views: NavigationView[]
  isMember: boolean
}

export function WorkspaceIndex(
  workspaceID: string,
  workspaces: NavigationWorkspace[]
): number | undefined {
  let idx = 0

  for (const ws of workspaces) {
    if (ws.id === workspaceID) {
      return idx
    }
    idx++

    for (const suggestion of ws.suggestingWorkspaces) {
      if (suggestion.id === workspaceID) {
        return idx
      }
      idx++
    }
  }

  return undefined
}

<template>
  <nav
    class="bg-gray-200 border-r border-gray-300 h-screen fixed space-y-1 h-screen flex flex-col z-20"
    aria-label="Sidebar"
  >
    <AppTitleBarSpacer v-slot="{ ipc }" pad-left="1rem" class="flex-none">
      <AppHistoryNavigationButtons :ipc="ipc" />
    </AppTitleBarSpacer>

    <div class="flex-shrink-0 border-b border-warmgray-300 h-16 z-20">
      <NavDropdown :user="user" @logout="$emit('logout')" />
    </div>

    <template v-if="false">
      <!-- Just a test to easily see that a page is tested in all sizes -->
      <div class="block sm:hidden">default</div>
      <div class="hidden sm:block md:hidden">sm</div>
      <div class="hidden md:block lg:hidden">md</div>
      <div class="hidden lg:block xl:hidden">lg</div>
      <div class="hidden xl:block">xl</div>
    </template>

    <div class="flex-1 mt-1 overflow-y-auto">
      <StackedMenuLoading v-if="isLoading" />
      <StackedMenuEmpty v-else-if="navigation.length === 0" class="mt-8" />
      <div v-for="(codebase, codebaseIdx) in navigation" v-else :key="codebase.id" class="relative">
        <OnboardingStep
          id="FindingYourCodebase"
          :dependencies="['CreatingANewWorkspace', 'FindingYourWorkspace']"
          :enabled="codebase.workspaces.some((workspace) => workspace.isCurrent)"
        >
          <template #title>Codebase</template>
          <template #description>
            To get back to the overview of your codebase, here's where you should look! On the
            codebase page you'll see the files at the most recently saved change. It's also where
            you'll find a way to invite team members to collaborate.
          </template>

          <div
            :class="[
              codebase.isCurrent
                ? 'bg-warmgray-50  hover:text-gray-900'
                : 'hover:bg-warmgray-100  hover:text-gray-900',

              'flex items-center pl-3 pr-1 py-2 text-sm font-medium rounded-md mx-1 justify-between transition whitespace-nowrap space-x-2 h-10 cursor-pointer text-gray-700 my-0.5',
            ]"
          >
            <router-link
              :to="{
                name: 'codebaseHome',
                params: { codebaseSlug: codebase.slug },
              }"
              :href="codebase.href"
              :aria-current="codebase.isCurrent ? 'page' : undefined"
              class="flex-1 overflow-hidden text-ellipsis flex-shrink-0"
            >
              <span>
                {{ codebase.name }}
              </span>
            </router-link>
            <OnboardingStep
              v-if="authenticated && codebase.isMember"
              id="CreatingANewWorkspace"
              :dependencies="['FindingYourWorkspace', 'LandingAChange']"
              :enabled="codebase.workspaces.some((ws) => ws.isCurrent)"
            >
              <template #title>More Workspaces</template>

              <template #description>
                Workspaces are cheap! You can make new ones on the fly, and move your local
                directory between them. Changes you make stay in the workspace until you share them!
                Make one to quickly fix a bug, or try out an idea without interrupting your flow.
              </template>

              <Tooltip :y-direction="codebaseIdx === 0 ? 'down' : 'up'" x-direction="left">
                <template #tooltip>New Workspace</template>
                <button
                  class="p-1 hover:bg-warmgray-200 rounded-md cursor-pointer transition flex-shrink-0"
                  @click="createWorkspaceHandler(codebase.slug, codebase.id)"
                >
                  <PlusSmIcon class="w-5 h-5 hover:text-gray-900" />
                </button>
              </Tooltip>
            </OnboardingStep>
          </div>
        </OnboardingStep>
        <div
          v-if="codebase.workspaces.length > 0"
          class="flex flex-col overflow-hidden text-ellipsis my-0.5"
        >
          <template v-for="(workspace, workspaceIdx) in codebase.workspaces" :key="workspace.id">
            <OnboardingStep
              id="FindingYourWorkspace"
              :enabled="codebaseIdx === 0 && workspaceIdx === 0"
            >
              <template #title>Workspaces</template>
              <template #description>
                <div v-if="workspace.currentView == null && workspace.isOwnedByUser">
                  Look over here! This is a workspace – a place where you can make changes to the
                  codebase. Connect your local directory to the workspace and start making changes
                  in your own editor!
                </div>
                <div v-else-if="workspace.currentView != null && workspace.isOwnedByUser">
                  Look over here! This is a workspace – a place where you can make change to the
                  codebase. Your local directory is connected to it, so you can start making changes
                  in your own editor!
                </div>
                <div v-else>
                  Look over here! This is a workspace – a place where changes are made to the
                  codebase. Create one for yourself and connect your local directory to it!
                </div>
              </template>

              <router-link
                :to="{
                  name: 'workspaceHome',
                  params: {
                    codebaseSlug: codebase.slug,
                    id: workspace.id,
                  },
                }"
                class="whitespace-nowrap text-gray-500 text-sm font-medium py-2 pl-2 pr-2 inline-flex items-center relative rounded-md my-0.5 mx-1 group"
                :class="[
                  workspace.isCurrent
                    ? 'bg-warmgray-50 hover:text-gray-900'
                    : 'hover:bg-warmgray-100 hover:text-gray-900',
                  !workspace.isCurrent && workspace.isUnread ? '!text-gray-900 !font-semibold' : '',
                ]"
              >
                <Avatar
                  size="5"
                  :author="workspace.author"
                  :class="[
                    workspace.isAuthorCoding ? 'ring-2 ring-green-400' : 'ring-1 ring-gray-300',
                  ]"
                  class="mr-2 flex-shrink-0"
                />
                <span :class="['text-ellipsis overflow-hidden flex-1']">
                  {{ workspace.name }}
                </span>
                <span
                  v-if="workspace.badgeCount > 0"
                  class="flex rounded-full bg-red-500 text-white h-7 -my-1 -mx-1 border justify-center whitespace-nowrap"
                  :class="[workspace.badgeCount < 100 ? 'w-7' : 'px-2']"
                >
                  <span class="text-xs leading-7 font-semibold">
                    {{ workspace.badgeCount }}
                  </span>
                </span>
                <!-- Make space for the view icon -->
                <span v-if="workspace.currentView != null" class="w-8"></span>
              </router-link>
              <template
                v-for="(
                  suggestingWorkspace, suggestingWorkspaceIdx
                ) in workspace.suggestingWorkspaces"
                :key="suggestingWorkspace.id"
              >
                <router-link
                  :to="{
                    name: 'workspaceHome',
                    params: {
                      codebaseSlug: codebase.slug,
                      id: suggestingWorkspace.id,
                    },
                  }"
                  class="whitespace-nowrap text-gray-500 text-sm font-medium py-2 pl-2 pr-2 inline-flex items-center relative rounded-md my-0.5 mx-1 ml-8 group"
                  :class="[
                    suggestingWorkspace.isCurrent
                      ? 'bg-warmgray-50 hover:text-gray-900'
                      : 'hover:bg-warmgray-100 hover:text-gray-900',
                  ]"
                >
                  <span :class="['text-ellipsis overflow-hidden flex-1']">
                    Suggestion {{ suggestingWorkspaceIdx + 1 }}
                  </span>
                  <!-- Make space for the view icon -->
                  <span v-if="workspace.currentView != null" class="w-8"></span>
                </router-link>
              </template>
            </OnboardingStep>
          </template>
        </div>
        <template v-for="view in codebase.views" :key="view.id">
          <div class="absolute top-0 left-0 h-0 w-full">
            <!-- 36px = height of workspace list item -->
            <!-- 40px = height of codebase list item -->

            <div
              v-if="view.isConnectedToWorkspace && view.workspaceIndex !== undefined"
              class="absolute transition transition-transform ease-in-out duration-300 h-7 top-0 right-0"
              :style="{
                transform: `translateY(calc(40px + 0.125rem + ${view.workspaceIndex} * (36px + 0.25rem) + (36px + 0.25rem) / 2)) translateY(-50%)`,
              }"
            >
              <div class="h-7 absolute right-0 group flex mx-2 hover:pl-4">
                <div class="flex-1 group-hover:hidden"></div>
                <div
                  class="transition transition-background group-hover:flex-1 backdrop-blur-sm bg-gray-50 group-hover:bg-gray-100 p-1 border h-7 w-max group-hover:w-[15rem] rounded-full flex items-center justify-center overflow-hidden shadow-lg"
                >
                  <div
                    class="flex-1 whitespace-nowrap text-xs px-1 overflow-hidden text-ellipsis hidden group-hover:block"
                  >
                    <span class="font-medium">
                      {{ ellipsisShortPath(view.data.shortMountPath) }}
                    </span>

                    <span class="text-gray-500"> on {{ view.data.mountHostname }} </span>
                  </div>

                  <div class="flex-none">
                    <ViewStatusIndicator compact :view="view.data" />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>

    <div
      v-if="authenticated && data"
      class="flex-none px-3 py-2 flex items-center border-t border-warmgray-300 space-x-4 justify-between"
    >
      <router-link :to="{ name: 'user' }">
        <Avatar :author="user" size="10" class="border border-gray-300" />
      </router-link>

      <div class="flex space-x-2">
        <NotificationIcon
          v-if="user"
          :features="features"
          :user="user"
          class="justify-center inline-flex"
        />
        <router-link
          :to="{ name: 'user' }"
          class="p-2 hover:bg-warmgray-300 text-gray-400 hover:text-gray-700 transition rounded-md"
        >
          <CogIcon class="h-5 w-5" />
        </router-link>
      </div>
    </div>
  </nav>
</template>

<script lang="ts">
import { CogIcon, PlusSmIcon } from '@heroicons/vue/solid'
import { gql, useQuery } from '@urql/vue'
import { defineComponent, onUnmounted, PropType, ref } from 'vue'
import { useRoute } from 'vue-router'
import { IdFromSlug, Slug } from '../../slug'
import Avatar from '../shared/Avatar.vue'
import NavDropdown from '../NavDropdown.vue'
import NotificationIcon from '../notification/Icon.vue'
import StackedMenuEmpty from './StackedMenuEmpty.vue'
import StackedMenuLoading from './StackedMenuLoading.vue'
import { useUpdatedCodebase } from '../../subscriptions/useUpdatedCodebase'
import { useUpdatedWorkspaceActivity } from '../../subscriptions/useUpdatedWorkspaceActivity'
import { useCreateWorkspace } from '../../mutations/useCreateWorkspace'
import ViewStatusIndicator, { VIEW_STATUS_INDICATOR } from '../ViewStatusIndicator.vue'
import { useOpenWorkspaceOnView } from '../../mutations/useOpenWorkspaceOnView'
import OnboardingStep from '../onboarding/OnboardingStep.vue'
import Tooltip from '../shared/Tooltip.vue'
import {
  CodebaseFragment,
  PresenceFragment as WorkspacePresence,
  StackedMenu_ViewFragment as View,
  StackedMenu_WorkspaceActivityFragment as WorkspaceActivity,
  StackedMenu_WorkspaceReviewFragment as WorkspaceReview,
  WorkspaceFragment,
} from './__generated__/StackedMenu'
import { ReviewGrade, User, WorkspacePresenceState, Feature } from '../../__generated__/types'
import { useUpdatedReviews } from '../../subscriptions/useUpdatedReviews'
import {
  NavigationCodebase,
  NavigationSuggestingWorkspace,
  NavigationView,
  NavigationWorkspace,
  WorkspaceIndex,
} from './MenuHelper'
import AppTitleBarSpacer from '../AppTitleBarSpacer.vue'
import AppHistoryNavigationButtons from '../AppHistoryNavigationButtons.vue'
import { useUpdatedViews } from '../../subscriptions/useUpdatedViews'

const WORKSPACE_ACTIVITY_FRAGMENT = gql`
  fragment StackedMenu_WorkspaceActivity on WorkspaceActivity {
    id
    isRead
    ... on WorkspaceRequestedReviewActivity {
      id
    }
    ... on WorkspaceCommentActivity {
      comment {
        id
        message
      }
    }
  }
`

const PRESENSE_FRAGMENT = gql`
  fragment Presence on WorkspacePresence {
    id
    author {
      id
      name
      avatarUrl
    }
    state
    lastActiveAt
  }
`

const WORKSPACE_REVIEW_FRAGMENT = gql`
  fragment StackedMenu_WorkspaceReview on Review {
    id
    author {
      id
    }
    dismissedAt
    isReplaced
    grade
  }
`

const WORKSPACE_AUTHOR_FRAGMENT = gql`
  fragment StackedMenu_WorkspaceAuthor on Author {
    id
    name
    avatarUrl
  }
`

const WORKSPACE_FRAGMENT = gql`
  fragment Workspace on Workspace {
    id
    name
    lastActivityAt
    createdAt

    author {
      ...StackedMenu_WorkspaceAuthor
    }

    presence {
      ...Presence
    }

    suggestion {
      id
      for {
        id
        name
      }
      createdAt
    }

    activity {
      ...StackedMenu_WorkspaceActivity
    }

    reviews {
      ...StackedMenu_WorkspaceReview
    }
  }
  ${PRESENSE_FRAGMENT}
  ${WORKSPACE_ACTIVITY_FRAGMENT}
  ${WORKSPACE_REVIEW_FRAGMENT}
  ${WORKSPACE_AUTHOR_FRAGMENT}
`

const VIEW_FRAGMENT = gql`
  fragment StackedMenu_View on View {
    id
    shortMountPath
    mountHostname
    workspace {
      id
    }
    ...ViewStatusIndicator
  }
`

const MEMBER_FRAGMENT = gql`
  fragment StackedMenu_Member on Author {
    id
  }
`

const CODEBASE_FRAGMENT = gql`
  fragment Codebase on Codebase {
    id
    name
    shortID
    lastUpdatedAt
    archivedAt
    createdAt
    isReady

    workspaces {
      ...Workspace
    }

    views {
      ...StackedMenu_View
    }

    members {
      ...StackedMenu_Member
    }
  }

  ${VIEW_STATUS_INDICATOR}
  ${WORKSPACE_FRAGMENT}
  ${VIEW_FRAGMENT}
  ${MEMBER_FRAGMENT}
`

const nonArchived = (codebase: CodebaseFragment) => !codebase.archivedAt

const onlyReady = (codebase: CodebaseFragment) => codebase.isReady

const codebaseByMembership = function (
  user: User | undefined
): (a: CodebaseFragment, b: CodebaseFragment) => number {
  if (!user) return () => 0
  return (a, b): number => {
    const memberOfA = a.members.some(({ id }) => id === user.id)
    const memberOfB = b.members.some(({ id }) => id === user.id)
    if (memberOfA && !memberOfB) {
      return 1
    } else if (!memberOfA && memberOfB) {
      return -1
    } else {
      return 0
    }
  }
}

const codebaseByLastUpdated = (a: CodebaseFragment, b: CodebaseFragment) => {
  let av = Math.max(a.lastUpdatedAt || 0, a.createdAt)
  let bv = Math.max(b.lastUpdatedAt || 0, b.createdAt)

  if (!av && !bv) return 0
  if (!bv) return -1
  if (!av) return 1
  return bv - av
}

const ownedBy = function (user: User | undefined): (ws: WorkspaceFragment) => boolean {
  return (ws: WorkspaceFragment) => ws.author.id === user?.id
}

const nonSuggestingWorkspaces = (ws: WorkspaceFragment) => !ws.suggestion
const suggestingWorkspaces = (ws: WorkspaceFragment) => !nonSuggestingWorkspaces(ws)

const workspaceByLastUpdated = (a: WorkspaceFragment, b: WorkspaceFragment) =>
  b.lastActivityAt - a.lastActivityAt

const isPresenseCoding = (p: WorkspacePresence) => {
  return p.state === WorkspacePresenceState.Coding
}

const isPresenceNotStale = (p: WorkspacePresence) => {
  const now = Math.floor(new Date().getTime() / 1000)
  return p.lastActiveAt > now - 600
}

const isWorkspaceAuthorCoding = (ws: WorkspaceFragment) => {
  return (
    ws.presence
      .filter((presence) => presence.author.id === ws.author.id)
      .filter(isPresenseCoding)
      .filter(isPresenceNotStale).length > 0
  )
}

const isActivityUnread = (activity: WorkspaceActivity): boolean => activity.isRead === false

const isReviewRequested = (review: WorkspaceReview) => {
  return !review.dismissedAt && !review.isReplaced && review.grade === ReviewGrade.Requested
}

const viewByWorkspaceId = function (workspaceId: string): (view: View) => boolean {
  return (view: View) => view.workspace?.id === workspaceId
}

const reviewByUserId = function (userId: string): (review: WorkspaceReview) => boolean {
  return (review: WorkspaceReview) => review.author.id === userId
}

const hasMention = function (userId: string): (activity: WorkspaceActivity) => boolean {
  return (activity: WorkspaceActivity): boolean => {
    if (activity.__typename !== 'WorkspaceCommentActivity') return false
    if (!activity.comment) return false
    if (activity.comment.message.indexOf(`@${userId}`) === -1) return false
    return true
  }
}

export default defineComponent({
  components: {
    Tooltip,
    OnboardingStep,
    ViewStatusIndicator,
    StackedMenuEmpty,
    NavDropdown,
    PlusSmIcon,
    Avatar,
    CogIcon,
    NotificationIcon,
    StackedMenuLoading,
    AppHistoryNavigationButtons,
    AppTitleBarSpacer,
  },

  props: {
    user: {
      type: Object as PropType<User>,
    },
    features: {
      type: Array as PropType<Feature[]>,
      required: true,
    },
  },

  emits: ['logout'],

  setup() {
    const route = useRoute()
    const shortCodebaseID = ref(
      route.params.codebaseSlug ? IdFromSlug(route.params.codebaseSlug as string) : ''
    )

    const { data, executeQuery, fetching } = useQuery({
      query: gql`
        query StackedMenu($shortCodebaseId: ID!) {
          codebase(shortID: $shortCodebaseId) {
            ...Codebase
          }

          codebases {
            ...Codebase
          }
        }
        ${CODEBASE_FRAGMENT}
      `,

      variables: {
        shortCodebaseId: shortCodebaseID,
      },
      requestPolicy: 'cache-and-network',
    })

    useUpdatedWorkspaceActivity()
    useUpdatedCodebase()
    useUpdatedReviews()
    useUpdatedViews()

    // Reload data every 15s
    const reload = () => {
      executeQuery({ requestPolicy: 'network-only' })
    }
    const interval = setInterval(reload, 15000)
    onUnmounted(() => {
      clearInterval(interval)
    })

    const createWorkspaceResult = useCreateWorkspace()
    const openWorkspaceOnView = useOpenWorkspaceOnView()
    return {
      data,
      shortCodebaseID,

      fetching,

      route,

      openWorkspaceOnView,
      createWorkspace(codebaseID: string) {
        return createWorkspaceResult({
          codebaseID,
        })
      },
    }
  },

  computed: {
    isLoading(): boolean {
      if (this.data) return false
      return this.fetching
    },
    authenticated(): boolean {
      return !!this.user
    },
    codebases(): CodebaseFragment[] {
      if (!this.data) return []
      const codebases = this.data.codebase
        ? [
            ...this.data.codebases.filter((cb: CodebaseFragment) => cb.id != this.data.codebase.id),
            this.data.codebase,
          ]
        : this.data.codebases
      return codebases
        .filter(nonArchived)
        .filter(onlyReady)
        .sort(codebaseByLastUpdated)
        .sort(codebaseByMembership(this.user))
    },
    navigation(): NavigationCodebase[] {
      return this.codebases.map((cb: CodebaseFragment): NavigationCodebase => {
        const shortCodebaseID = this.route.params.codebaseSlug
          ? IdFromSlug(this.route.params.codebaseSlug as string)
          : ''
        const isMember = cb.members.some(({ id }) => id === this.user?.id)

        const suggestions = cb.workspaces
          .filter(suggestingWorkspaces)
          .filter(ownedBy(this.user))
          .reduce((acc, current) => {
            if (!current.suggestion) return acc

            let suggesting = acc.get(current.suggestion.for.id)
            if (!suggesting) suggesting = []
            suggesting.push({
              id: current.id,
              isCurrent: current.id === this.route.params.id,
              createdAt: current.suggestion.createdAt,
            })
            suggesting = suggesting.sort((a, b) => a.createdAt - b.createdAt)
            acc.set(current.suggestion.for.id, suggesting)
            return acc
          }, new Map<string, NavigationSuggestingWorkspace[]>())

        const workspaces = cb.workspaces
          .sort(workspaceByLastUpdated)
          .filter(nonSuggestingWorkspaces)
          .map((ws: WorkspaceFragment): NavigationWorkspace => {
            const requestedReviews = this.user
              ? ws.reviews.filter(reviewByUserId(this.user.id)).filter(isReviewRequested).length
              : 0

            const unseenMentions = this.user
              ? ws.activity.filter(isActivityUnread).filter(hasMention(this.user.id)).length
              : 0

            return {
              id: ws.id,
              name: ws.name,
              isCurrent: shortCodebaseID === cb.shortID && ws.id === this.route.params.id,
              author: ws.author,
              isUnread: isMember ? ws.activity.some(isActivityUnread) : false,
              badgeCount: unseenMentions + requestedReviews,
              presence: ws.presence,
              isAuthorCoding: isWorkspaceAuthorCoding(ws),
              currentView: cb.views.find(viewByWorkspaceId(ws.id)),
              isOwnedByUser: !!this.user && ws.author.id === this.user.id,
              suggestingWorkspaces: suggestions.get(ws.id) || [],
            }
          })

        const views = cb.views.map((view): NavigationView => {
          const workspaceID = view.workspace?.id
          const workspaceIndex = workspaceID ? WorkspaceIndex(workspaceID, workspaces) : undefined
          const isConnectedToWorkspace = workspaceIndex !== undefined

          return {
            id: view.id,
            isConnectedToWorkspace,
            workspaceId: workspaceID,
            workspaceIndex,
            data: view,
          }
        })

        return {
          id: cb.id,
          name: cb.name,
          slug: Slug(cb.name, cb.shortID),
          isCurrent: shortCodebaseID === cb.shortID && !workspaces.some((ws) => ws.isCurrent),
          isMember: isMember,
          workspaces: workspaces,
          views: views,
        }
      })
    },
  },

  watch: {
    $route: function (newRoute) {
      this.shortCodebaseID = newRoute.params.codebaseSlug
        ? IdFromSlug(newRoute.params.codebaseSlug as string)
        : ''
    },
  },

  methods: {
    async createWorkspaceHandler(codebaseSlug: string, codebaseId: string) {
      const res = await this.createWorkspace(codebaseId)

      this.$router.push({
        name: 'workspaceHome',
        params: { codebaseSlug: codebaseSlug, id: res.createWorkspace.id },
      })
    },

    ellipsisShortPath(path: string): string {
      const slicedPath = path.slice(-50)
      if (slicedPath === path) {
        return path
      } else {
        return `…${slicedPath.replace(/^[^/]+/, '')}`
      }
    },
  },
})
</script>

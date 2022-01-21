import {
  createMemoryHistory,
  createRouter,
  createWebHashHistory,
  createWebHistory,
  RouteRecordRaw,
  RouterOptions,
} from 'vue-router'

// Public website imports
import Index from './pages/Index.vue' // Pre loaded

export const routes: RouteRecordRaw[] = [
  {
    path: '/',
    components: { default: Index },
    name: 'index',
    alias: '/index',
    meta: { nonApp: true },
  },
  {
    path: '/codebases',
    component: () => import('./pages/CodebaseOverview.vue'),
    name: 'codebaseOverview',
  },
  {
    path: '/org/:organizationSlug',
    component: () => import('./pages/organization/CodebaseListPage.vue'),
    name: 'organizationListCodebases',
  },
  { path: '/new', component: () => import('./pages/CreateEmpty.vue'), name: 'codebaseCreate' },
  {
    path: '/org/new',
    component: () => import('./pages/organization/CreateOrganizationPage.vue'),
    name: 'organizationCreate',
  },
  {
    path: '/org/:organizationSlug/settings',
    component: () => import('./pages/organization/OrganizationSettingsPage.vue'),
    name: 'organizationSettings',
  },
  {
    path: '/org/:organizationSlug/new',
    component: () => import('./pages/organization/CreateOrganizationCodebasePage.vue'),
    name: 'organizationCreateCodebase',
  },
  {
    path: '/org/:organizationSlug/subscriptions/new',
    component: () => import('./pages/organization/CreateSubscriptionPage.vue'),
    name: 'organizationCreateSubscription',
  },
  {
    path: '/setup-github',
    component: () => import('./pages/SetupGithub.vue'),
    name: 'setupGithub',
    meta: { selfContainedLayout: true },
  },
  {
    path: '/:codebaseSlug/settings',
    component: () => import('./pages/settings/Settings.vue'),
    name: 'codebaseSettings',
  },
  {
    path: '/:codebaseSlug/settings/team',
    component: () => import('./pages/settings/team/SettingsTeam.vue'),
    name: 'codebaseSettingsTeam',
  },
  {
    path: '/:codebaseSlug/settings/acl',
    component: () => import('./pages/settings/acl/SettingsACL.vue'),
    name: 'codebaseSettingsAcls',
  },
  {
    path: '/:codebaseSlug/settings/workspaces',
    component: () => import('./pages/settings/workspaces/SettingsWorkspaces.vue'),
    name: 'codebaseSettingsWorkspaces',
  },
  {
    path: '/:codebaseSlug/settings/integrations',
    component: () => import('./pages/settings/integrations/ListIntegrations.vue'),
    name: 'codebaseSettingsIntegrations',
  },
  {
    path: '/:codebaseSlug/settings/add/buildkite',
    component: () => import('./pages/settings/integrations/add/buildkite/Buildkite.vue'),
    name: 'codebaseSettingsAddBuildkite',
  },
  {
    path: '/:codebaseSlug/settings/edit/buildkite/:integrationId',
    component: () => import('./pages/settings/integrations/add/buildkite/Buildkite.vue'),
    name: 'codebaseSettingsEditBuildkite',
  },
  {
    path: '/:codebaseSlug/changes/:selectedChangeID?',
    component: () => import('./pages/Changelog.vue'),
    name: 'codebaseChangelog',
  },
  {
    path: '/auth/:email?',
    component: () => import('./components/auth/UserAuth.vue'),
    name: 'authWithEmail',
    meta: { nonApp: true, hideNavigation: true, isAuth: true },
  },
  {
    path: '/auth',
    component: () => import('./components/auth/UserAuth.vue'),
    name: 'auth',
    meta: { nonApp: true, hideNavigation: true, isAuth: true },
  },
  {
    path: '/signup',
    component: () => import('./components/auth/UserAuth.vue'),
    name: 'signup',
    meta: { nonApp: true, hideNavigation: true, isAuth: true },
  },
  {
    path: '/login',
    component: () => import('./components/auth/UserAuth.vue'),
    name: 'login',
    meta: { nonApp: true, hideNavigation: true, isAuth: true },
  },
  {
    path: '/download',
    component: () => import('./pages/download/DownloadPage.vue'),
    name: 'download',
    meta: { nonApp: true },
  },
  {
    path: '/get-started/yc',
    component: () => import('./components/auth/UserAuth.vue'),
    name: 'getStartedYC',
    meta: { nonApp: true, hideNavigation: true, isAuth: true },
  },
  {
    path: '/get-started/github',
    component: () => import('./components/auth/UserAuth.vue'),
    name: 'getStartedGitHub',
    meta: { nonApp: true, hideNavigation: true, isAuth: true },
  },
  { path: '/user', component: () => import('./pages/User.vue'), name: 'user' },
  {
    path: '/install/:codebaseSlug?',
    component: () => import('./pages/install/InstallClient.vue'),
    name: 'installClient',
    meta: { neverElectron: true },
  },
  {
    path: '/install/token',
    component: () => import('./pages/install/InstallToken.vue'),
    name: 'installToken',
    meta: { neverElectron: true },
  },
  {
    path: '/about',
    component: () => import('./pages/about/About.vue'),
    name: 'about',
    meta: { nonApp: true },
  },
  {
    path: '/contact',
    component: () => import('./pages/about/Contact.vue'),
    name: 'contact',
    meta: { nonApp: true },
  },
  {
    path: '/privacy',
    component: () => import('./pages/about/Privacy.vue'),
    name: 'privacy',
    meta: { nonApp: true },
  },
  {
    path: '/terms-of-service',
    component: () => import('./pages/about/ToS.vue'),
    name: 'termsOfService',
    meta: { nonApp: true },
  },
  {
    path: '/pricing',
    component: () => import('./pages/static/Pricing.vue'),
    name: 'pricing',
    meta: { nonApp: true },
  },
  {
    path: '/security',
    component: () => import('./pages/static/Security.vue'),
    name: 'resourcesSecurity',
    meta: {
      nonApp: true,
      documentation: { title: 'Security', group: 'Documentation' },
    },
  },
  {
    path: '/docs',
    component: () => import('./pages/static/Docs.vue'),
    name: 'resourcesDocs',
    meta: { nonApp: true },
  },
  {
    path: '/kth',
    component: () => import('./components/campaigns/Kth.vue'),
    name: 'kth',
    meta: { nonApp: true },
  },
  {
    path: '/aoc',
    redirect: '/advent-of-code-2021-uoeIDQk',
    name: 'advent-of-code',
    meta: { nonApp: true },
  },
  {
    path: '/features/access-control',
    component: () => import('./pages/static/access-control/AccessControl.vue'),
    name: 'featuresAccessControl',
    meta: {
      nonApp: true,
    },
  },
  {
    path: '/features/instant-integration',
    component: () => import('./pages/static/instant-integration/InstantIntegration.vue'),
    name: 'featuresInstantIntegration',
    meta: {
      nonApp: true,
    },
  },
  {
    path: '/quickstart', // TODO: Move to /docs/quickstart
    component: () => import('./pages/static/Quickstart.vue'),
    name: 'quickstart',
    meta: {
      nonApp: true,
      documentation: {
        title: 'Get up and running with Sturdy',
        group: 'Quickstart',
      },
    },
  },
  {
    path: '/syncing', // TODO: Move to /docs/syncing
    component: () => import('./pages/static/SharingAndSyncing.vue'),
    name: 'syncing',
    meta: {
      nonApp: true,
      documentation: {
        title: 'Staying up to date with changes done by others',
        group: 'Quickstart',
      },
    },
  },
  {
    path: '/features/instant-switching',
    component: () => import('./pages/static/WorkspaceNavigation.vue'),
    name: 'featuresWorkspaceNavigation',
    meta: {
      nonApp: true,
      documentation: {
        title: 'Instant Workspace switching',
        group: 'Solutions',
      },
    },
  },
  {
    path: '/features/integrations',
    component: () => import('./pages/static/Integrations.vue'),
    name: 'featuresIntegrations',
    meta: {
      nonApp: true,
      documentation: {
        title: 'Integrate Sturdy with other tools',
        group: 'Solutions',
      },
    },
  },
  {
    path: '/features/live',
    component: () => import('./pages/static/LiveFeedback.vue'),
    name: 'featuresLiveFeedback',
    meta: {
      nonApp: true,
      documentation: { title: 'Live Feedback', group: 'Solutions' },
    },
  },
  {
    path: '/features/workflow',
    component: () => import('./pages/static/Workflow.vue'),
    name: 'featuresWorkflow',
    meta: { nonApp: true },
  },
  {
    path: '/features/conflicts',
    component: () => import('./pages/static/Conflicts.vue'),
    name: 'featuresConflicts',
    meta: {
      nonApp: true,
      documentation: {
        title: 'Conflict Resolution on Sturdy',
        group: 'Solutions',
      },
    },
  },
  {
    path: '/features/large-files',
    component: () => import('./pages/static/LargeFiles.vue'),
    name: 'featuresLargeFiles',
    meta: {
      nonApp: true,
      documentation: { title: 'Large Files', group: 'Solutions' },
    },
  },
  {
    path: '/docs/cli',
    component: () => import('./pages/static/HelpSturdyCommandLine.vue'),
    name: 'resourcesHelpSturdyCommandLine',
    meta: {
      nonApp: true,
      documentation: {
        title: 'The Sturdy command line application',
        group: 'Help',
      },
    },
  },
  {
    path: '/api',
    component: () => import('./pages/static/API.vue'),
    name: 'resourcesApi',
    meta: {
      nonApp: true,
      documentation: { title: 'GraphQL API', group: 'Documentation' },
    },
  },
  {
    path: '/docs/access-control',
    component: () => import('./pages/static/docs/access-control/AccessControl.vue'),
    name: 'docsAccessControl',
    meta: {
      nonApp: true,
      documentation: { title: 'Access Control', group: 'Documentation' },
    },
  },
  {
    path: '/docs/suggestions',
    component: () => import('./pages/static/docs/suggestions/Suggestions.vue'),
    name: 'docsSuggestions',
    meta: {
      nonApp: true,
      documentation: { title: 'Suggestions', group: 'Documentation' },
    },
  },
  {
    path: '/docs/sturdy-for-git-users',
    component: () => import('./pages/static/docs/terminology/Terminology.vue'),
    name: 'docsSturdyForGitUsers',
    meta: {
      nonApp: true,
      documentation: { title: 'Sturdy for Git users', group: 'Documentation' },
    },
  },
  {
    path: '/features/migrate-from-github',
    component: () => import('./pages/static/MigrateFromGitHub.vue'),
    name: 'resourcesMigrateFromGitHub',
    meta: {
      nonApp: true,
      documentation: { title: 'Migrate from GitHub', group: 'Documentation' },
    },
  },
  {
    path: '/docs/continuous-integration',
    component: () => import('./pages/static/continuous-integration/ContinuousIntegration.vue'),
    name: 'docsCICD',
    meta: {
      nonApp: true,
      documentation: { title: 'CI/CD', group: 'Documentation' },
    },
  },
  {
    path: '/blog/2021-12-17-graphql-componentized-uis',
    name: 'graphql-componentized-uis',
    component: () => import('./pages/blog/graphql-componentized-uis/Post.vue'),
    meta: {
      nonApp: true,
      blog: { title: '#013 - GraphQL & Componentized UIs' },
    },
  },
  {
    path: '/blog/2021-12-07-launching-the-sturdy-app',
    name: 'blog2021-12-07',
    component: () => import('./pages/blog/2021-12-07/Post.vue'),
    meta: {
      nonApp: true,
      blog: { title: '#012 - Launching the Sturdy App!' },
    },
  },
  {
    path: '/blog/2021-11-29-scaling-teams',
    name: 'scaling-teams',
    component: () => import('./pages/blog/scaling-teams/Post.vue'),
    meta: {
      nonApp: true,
      blog: { title: '#011 - Scaling teams as parallel computing systems' },
    },
  },
  {
    path: '/blog/2021-11-22-sturdy-the-app-is-coming',
    name: 'blog-2021-11-22',
    component: () => import('./pages/blog/2021-11-22/Post.vue'),
    meta: { nonApp: true, blog: { title: '#010 - Sturdy-the-app is coming!' } },
  },
  {
    path: '/blog/2021-09-29-acls-and-a-fresh-hot-look',
    name: 'blogRedesign',
    component: () => import('./pages/blog/redesign/Redesign.vue'),
    meta: { nonApp: true, blog: { title: '#009 - ACLs and a fresh hot look!' } },
  },
  {
    path: '/blog/2021-09-09-large-files',
    name: 'blogLargeFiles',
    component: () => import('./pages/blog/large-files/LargeFiles.vue'),
    meta: { nonApp: true, blog: { title: '#008 - Large Files' } },
  },
  {
    path: '/blog/2021-08-18-unbreaking-code-collaboration',
    name: 'blogVision',
    component: () => import('./pages/blog/vision/Vision.vue'),
    meta: {
      nonApp: true,
      blog: { title: '#007 - Unbreaking code collaboration' },
    },
  },
  {
    path: '/blog/2021-08-12-signup-is-open',
    name: 'blogSignupIsOpen',
    component: () => import('./pages/blog/signup-is-open/SignupIsOpen.vue'),
    meta: { nonApp: true, blog: { title: '#006 - Sturdy is here' } },
  },
  {
    path: '/blog/2021-06-10-humane-code-review',
    name: 'blogHumaneCodeReview',
    component: () => import('./pages/blog/humane-code-review/HumaneCodeReview.vue'),
    meta: { nonApp: true, blog: { title: '#005 - Humane Code Review' } },
  },
  {
    path: '/blog/2021-05-06-importing-from-git',
    name: 'blogImportingFromGit',
    component: () => import('./pages/blog/import-from-git/ImportFromGit.vue'),
    meta: { nonApp: true, blog: { title: '#004 - Importing from Git' } },
  },
  {
    path: '/blog/2021-04-16-share-now',
    name: 'blogShareNow',
    component: () => import('./pages/blog/share-now/ShareNow.vue'),
    meta: { nonApp: true, blog: { title: '#003 - Share Now!' } },
  },
  {
    path: '/blog/2021-04-01-restore-to-any-point-in-time',
    name: 'blogRestoreToAnyPointInTime',
    component: () =>
      import('./pages/blog/restore-to-any-point-in-time/RestoreToAnyPointInTime.vue'),
    meta: {
      nonApp: true,
      blog: { title: '#002 - Restore to any point in time' },
    },
  },
  {
    path: '/blog/2021-03-24-yc-w21-demo-day',
    name: 'blogDemoDay',
    component: () => import('./pages/blog/yc-w21-demo-day/DemoDay.vue'),
    meta: { nonApp: true, blog: { title: '#001 - YC W21 Demo Day!' } },
  },
  {
    path: '/blog/2021-03-18-this-week-at-sturdy',
    name: 'blogThisWeekAtSturdy',
    component: () => import('./pages/blog/first-post/ThisWeekAtSturdy.vue'),
    meta: { nonApp: true, blog: { title: '#000 - This Week at Sturdy' } },
  },
  {
    path: '/blog',
    name: 'blog',
    component: () => import('./pages/blog/Blog.vue'),
    meta: { nonApp: true },
  },
  {
    path: '/careers',
    component: () => import('./pages/careers/Careers.vue'),
    name: 'careers',
    meta: { nonApp: true, selfContainedLayout: true, neverElectron: true },
  },
  {
    path: '/careers/founding-backend-engineer',
    component: () => import('./pages/careers/FoundingBackendEngineer.vue'),
    name: 'careersFoundingBackendEngineer',
    meta: { nonApp: true, selfContainedLayout: true, neverElectron: true },
  },
  {
    path: '/careers/founding-frontend-engineer',
    component: () => import('./pages/careers/FoundingFrontendEngineer.vue'),
    name: 'careersFoundingFrontendEngineer',
    meta: { nonApp: true, selfContainedLayout: true, neverElectron: true },
  },
  {
    path: '/v2/docs',
    component: () => import('./pages/docs/DocsRoot.vue'),
    name: 'v2DocsRoot',
    meta: { nonApp: true, selfContainedLayout: true, neverElectron: true },
  },
  {
    path: '/v2/docs/product-intro',
    component: () => import('./pages/docs/SturdyProductIntro.vue'),
    name: 'v2DocsProductIntro',
    meta: { nonApp: true, selfContainedLayout: true, neverElectron: true },
  },
  {
    path: '/v2/docs/how-sturdy-augments-git',
    component: () => import('./pages/docs/HowSturdyAugmentsGit.vue'),
    name: 'v2DocsHowSturdyAugmentsGit',
    meta: { nonApp: true, selfContainedLayout: true, neverElectron: true },
  },
  {
    path: '/v2/docs/working-in-the-open',
    component: () => import('./pages/docs/WorkingInTheOpen.vue'),
    name: 'v2DocsWorkingInTheOpen',
    meta: { nonApp: true, selfContainedLayout: true, neverElectron: true },
  },
  {
    path: '/v2/docs/how-to-ship-software-to-production',
    component: () => import('./pages/docs/HowToShipSoftwareToProduction.vue'),
    name: 'v2DocsHotToShipSoftwareToProduction',
    meta: { nonApp: true, selfContainedLayout: true, neverElectron: true },
  },
  {
    path: '/v2/docs/how-to-collaborate-with-others',
    component: () => import('./pages/docs/HowToCollaborateWithOthers.vue'),
    name: 'v2DocsHowToCollaborateWithOthers',
    meta: { nonApp: true, selfContainedLayout: true, neverElectron: true },
  },
  {
    path: '/v2/docs/how-to-edit-code',
    component: () => import('./pages/docs/HowToEditCode.vue'),
    name: 'v2DocsHowToEditCode',
    meta: { nonApp: true, selfContainedLayout: true, neverElectron: true },
  },
  {
    path: '/v2/docs/how-to-setup-sturdy-on-github',
    component: () => import('./pages/docs/HowToSetupSturdyOnGithub.vue'),
    name: 'v2DocsHowToSetupSturdyOnGitHub',
    meta: { nonApp: true, selfContainedLayout: true, neverElectron: true },
  },
  {
    path: '/v2/docs/how-to-switch-between-tasks',
    component: () => import('./pages/docs/HowToSwitchBetweenTasks.vue'),
    name: 'v2DocsHowToSwitchBetweenTasks',
    meta: { nonApp: true, selfContainedLayout: true, neverElectron: true },
  },
  {
    path: '/v2/docs/quickstart',
    component: () => import('./pages/docs/QuickStart.vue'),
    name: 'v2DocsQuickStart',
    meta: { nonApp: true, selfContainedLayout: true, neverElectron: true },
  },
  {
    path: '/v2/docs/using-sturdy',
    component: () => import('./pages/docs/UsingSturdy.vue'),
    name: 'v2DocsUsingSturdy',
    meta: { nonApp: true, selfContainedLayout: true, neverElectron: true },
  },
  {
    path: '/join/:code',
    component: () => import('./components/join/Join.vue'),
    name: 'join',
    meta: { nonApp: true, hideNavigation: true, isAuth: true },
  },
  {
    path: '/:codebaseSlug',
    component: () => import('./pages/CodebaseHome.vue'),
    name: 'codebaseHome',
  },

  {
    path: '/verify',
    name: 'verify',
    component: () => import('./components/emails/Verify.vue'),
    meta: { hideNavigation: true, isAuth: false },
  },
  {
    path: '/unsubscribe/:email',
    name: 'unsubscribe',
    component: () => import('./pages/newsletter/Unsubscribe.vue'),
    meta: { nonApp: true, hideNavigation: true, isAuth: false, skipPrerender: true },
  },
  {
    path: '/:codebaseSlug/:id',
    component: () => import('./pages/WorkspaceHome.vue'),
    name: 'workspaceHome',
  },
  {
    path: '/:codebaseSlug/browse/:path(.*)*',
    component: () => import('./pages/BrowseFile.vue'),
    name: 'browseFile',
  },
  {
    path: '/:codebaseSlug/browse',
    redirect: (to) => ({
      replace: true,
      path: to.fullPath + '/',
    }),
  },
]

const routerOpts: Partial<RouterOptions> = {
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else if (to.hash) {
      return { el: to.hash, scrollBehavior: 'smooth', top: 150 }
    } else {
      return { top: 0 }
    }
  },
}

if (import.meta.env.SSR) {
  routerOpts.history = createMemoryHistory()
} else if (import.meta.env.VITE_USE_HASH_HISTORY) {
  routerOpts.history = createWebHashHistory()
} else {
  routerOpts.history = createWebHistory()
}

export default createRouter(routerOpts as RouterOptions)

import { RouteRecordRaw } from 'vue-router'

export const RoutesSelfHosted: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: { name: 'login' },
  },
]

export const RoutesApps: RouteRecordRaw[] = [
  {
    path: '/home',
    component: () => import('./pages/HomePage.vue'),
    name: 'home',
    alias: '/codebases',
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
    path: '/org/:organizationSlug/settings/github',
    component: () => import('./pages/organization/OrganizationSetupGitHubPage.vue'),
    name: 'organizationSettingsGitHub',
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
    path: '/org/:organizationSlug/subscriptions',
    component: () => import('./pages/organization/ListSubscriptionPage.vue'),
    name: 'organizationListSubscription',
  },
  {
    path: '/org/:organizationSlug/installation',
    component: () => import('./pages/organization/ManageInstallationPage.vue'),
    name: 'organizationManageInstallation',
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
    path: '/:codebaseSlug/settings/developers',
    component: () => import('./pages/settings/developers/SettingsDevelopers.vue'),
    name: 'codebaseSettingsDevelopers',
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
    path: '/:codebaseSlug/changes',
    component: () => import('./pages/changes/List.vue'),
    name: 'codebaseChanges',
  },
  {
    path: '/:codebaseSlug/changes/:selectedChangeID',
    component: () => import('./pages/changes/Change.vue'),
    name: 'codebaseChange',
  },
  {
    path: '/auth/:email?',
    component: () => import('./pages/LoginRegister.vue'),
    name: 'authWithEmail',
    meta: { nonApp: true, hideNavigation: true, isAuth: true },
  },
  {
    path: '/auth',
    component: () => import('./pages/LoginRegister.vue'),
    name: 'auth',
    meta: { nonApp: true, hideNavigation: true, isAuth: true },
  },
  {
    path: '/signup',
    component: () => import('./pages/LoginRegister.vue'),
    name: 'signup',
    meta: { nonApp: true, hideNavigation: true, isAuth: true },
    props: {
      startWithSignUp: true,
    },
  },
  {
    path: '/login',
    component: () => import('./pages/LoginRegister.vue'),
    name: 'login',
    meta: { nonApp: true, hideNavigation: true, isAuth: true },
  },

  { path: '/user', component: () => import('./pages/User.vue'), name: 'user' },

  {
    path: '/install/:codebaseSlug?',
    component: () => import('./pages/install/InstallClient.vue'),
    name: 'installClient',
  },

  {
    path: '/install/token',
    component: () => import('./pages/install/InstallToken.vue'),
    name: 'installToken',
  },

  {
    path: '/privacy',
    component: () => import('./pages/docs/about/Privacy.vue'),
    name: 'privacy',
    meta: { nonApp: true, selfContainedLayout: true, neverElectron: true },
  },
  {
    path: '/terms-of-service',
    component: () => import('./pages/docs/about/ToS.vue'),
    name: 'termsOfService',
    meta: { nonApp: true, selfContainedLayout: true, neverElectron: true },
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
  {
    path: '/join/:code',
    component: () => import('./components/join/Join.vue'),
    name: 'join',
    meta: { nonApp: true, hideNavigation: true, isAuth: true },
  },
]

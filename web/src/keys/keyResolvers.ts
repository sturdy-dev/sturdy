import type { KeyingConfig } from '@urql/exchange-graphcache'

export const keyResolvers: KeyingConfig = {
  NotificationPreference: (data) => `${data.channel}/${data.type}`,
  WorkspaceWatcher: () => null,
  LicenseMessage: (data) => `${data.type}/${data.level}/${data.message}`,
}

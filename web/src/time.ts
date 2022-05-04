import { formatDistanceStrict, formatDistance } from 'date-fns'

export const getRelativeTimeStrict = (date: Date, baseDate = new Date()): string =>
  formatDistanceStrict(date, baseDate, { addSuffix: true })
export const getRelativeTime = (date: Date, baseDate = new Date()): string =>
  formatDistance(date, baseDate, { addSuffix: true })

export default {
  getRelativeTime,
}

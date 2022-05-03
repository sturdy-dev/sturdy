import { formatDistanceStrict } from 'date-fns'

export default {
  getRelativeTime: (date: Date, baseDate = new Date()): string =>
    formatDistanceStrict(date, baseDate, { addSuffix: true }),
}

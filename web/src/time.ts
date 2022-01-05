const units: { [key: string]: number } = {
  year: 24 * 60 * 60 * 1000 * 365,
  month: (24 * 60 * 60 * 1000 * 365) / 12,
  day: 24 * 60 * 60 * 1000,
  hour: 60 * 60 * 1000,
  minute: 60 * 1000,
  second: 1000,
}

const rtf = new Intl.RelativeTimeFormat('en', { numeric: 'auto' })

export default {
  getRelativeTime: function (d1: Date, d2 = new Date(), justNowMessage = 'just now'): string {
    const elapsed = d1.getTime() - d2.getTime()

    if (Math.abs(elapsed) < units.minute) {
      return justNowMessage
    }

    // "Math.abs" accounts for both "past" & "future" scenarios
    for (const u in units) {
      if (Math.abs(elapsed) > units[u] || u === 'second') {
        return rtf.format(Math.round(elapsed / units[u]), u as Intl.RelativeTimeFormatUnit)
      }
    }

    return ''
  },
}

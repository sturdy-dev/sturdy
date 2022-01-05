export const Slug = (humanName: string, shortID: string): string => {
  let safe = humanName.replace(/[^a-z0-9]/gi, '-').toLowerCase()
  safe = trimByChar(safe, '-')
  return safe + '-' + shortID
}

export const IdFromSlug = (slug: string): string => {
  const last = slug.lastIndexOf('-')
  return slug.substring(last + 1)
}

const trimByChar = (str: string, character: string): string => {
  const first = [...str].findIndex((char) => char !== character)
  const last = [...str].reverse().findIndex((char) => char !== character)
  return str.substring(first, str.length - last)
}

import { gql } from '@urql/vue'
import type { AuthorFragment } from './__generated__/AvatarHelper'
import { crc32 } from '@foxglove/crc'

export const AUTHOR = gql`
  fragment Author on Author {
    id
    name
    avatarUrl
  }
`

export const initials = function (author: AuthorFragment): string {
  if (author && author.name) {
    const words = author.name.split(' ')
    if (words.length === 1) {
      return words[0].substring(0, 1).toUpperCase()
    }

    return (
      words[0].substring(0, 1).toUpperCase() + words[words.length - 1].substring(0, 1).toUpperCase()
    )
  }
  return ''
}

const str2ab = function (str: string): ArrayBufferView {
  const array = new Uint8Array(str.length)
  for (let i = 0; i < str.length; i++) {
    array[i] = str.charCodeAt(i)
  }
  return array
}

export const initialsBgColor = function (author: AuthorFragment): string {
  if (!author) {
    return 'bg-gray-200'
  }

  const c = Math.abs(crc32(str2ab(author.id)))

  const colors = [
    'bg-red-100',
    'bg-yellow-100',
    'bg-green-100',
    'bg-blue-100',
    'bg-indigo-100',
    'bg-pink-100',
  ]

  return colors[c % colors.length]
}

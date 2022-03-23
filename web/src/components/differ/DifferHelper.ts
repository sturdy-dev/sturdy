import type { Hunk } from '../../__generated__/types'

export const getIndicesOf = function (
  searchStr: string,
  str: string,
  caseSensitive: boolean
): number[] {
  const searchStrLen = searchStr.length
  if (searchStrLen === 0) {
    return []
  }
  const indices = []
  if (!caseSensitive) {
    str = str.toLowerCase()
    searchStr = searchStr.toLowerCase()
  }
  let length = 0,
    started = false
  const lines = str.split('\n')
  for (const line of lines) {
    length += line.length + 1
    // Conditions to keep forward
    {
      if (!started) {
        if (line.startsWith('@@ ')) {
          started = true
        }
        continue
      } else if (line.includes('no newline at end of file')) {
        continue
      }
    }
    const index = line.indexOf(searchStr)
    if (index > -1) {
      const currentIndex = length - line.length < 0 ? 0 : length - line.length
      indices.push(currentIndex + index)
    }
  }

  return indices
}

export const searchMatches = function (
  searchResult: Map<string, number[]> | undefined,
  hunks: Hunk[]
): Set<string> {
  const res = new Set<string>()

  if (!searchResult || searchResult.size === 0) {
    return res
  }

  const endsAt = new Map<string, [[number, number]]>()

  for (const hunk of hunks) {
    let starts = 0,
      ends = 0,
      started = false
    const lines = hunk.patch.split('\n')
    for (const line of lines) {
      starts = ends
      ends += line.length + 1 // add trimmed newline

      // Conditions to keep forward
      {
        if (!started) {
          if (line.startsWith('@@ ')) {
            started = true
          }
          continue
        }

        if (line.includes('No newline at end of file')) {
          continue
        }
      }

      if (!endsAt.has(hunk.id)) {
        endsAt.set(hunk.id, [[starts, ends]])
      } else {
        const ref = endsAt.get(hunk.id)
        if (ref) {
          ref.push([starts, ends])
        }
      }
    }
  }

  for (const [hunkID, foundIndexes] of searchResult) {
    const rangeTuples = endsAt.get(hunkID)
    if (!rangeTuples) {
      continue
    } else if (rangeTuples.length <= 0) {
      continue
    }
    for (const foundIndex of foundIndexes) {
      for (let i = 0; i < rangeTuples.length; i++) {
        const tuple = rangeTuples[i],
          starts = tuple[0],
          ends = tuple[1]
        if (foundIndex < starts) {
          break
        } else if (foundIndex < ends) {
          res.add(hunkID + '-' + i)
          break
        }
      }
    }
  }
  return res
}

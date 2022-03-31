export interface SearchableHunk {
  hunkID: string
  patch: string
}

export const getIndicesOf = function (
  searchStr: string,
  str: string,
  caseSensitive: boolean
): number[] {
  if (searchStr.length === 0) {
    return []
  }

  if (!caseSensitive) {
    str = str.toLowerCase()
    searchStr = searchStr.toLowerCase()
  }

  const indices = []

  let length = 0
  let started = false

  const lines = str.split('\n')

  for (const line of lines) {
    length += line.length + 1

    if (!started) {
      if (line.startsWith('@@ ')) {
        started = true
      }
      continue
    }

    if (line.includes('no newline at end of file')) {
      continue
    }

    const index = line.indexOf(searchStr)

    if (index > -1) {
      const currentIndex = length - line.length < 0 ? 0 : length - line.length
      indices.push(currentIndex + index)
    }
  }

  return indices
}

type endsVal = {
  starts: number
  ends: number
  rowIndex: number
  blockId: number
}

export const searchMatches = function (
  searchResult: Map<string, number[]> | undefined,
  hunks: SearchableHunk[]
): Set<string> {
  const res = new Set<string>()

  if (!searchResult || searchResult.size === 0) {
    return res
  }

  const endsAt = new Map<string, [endsVal]>()

  for (const hunk of hunks) {
    let starts = 0
    let ends = 0
    let blockId = 0
    let rowIndex = 0
    let started = false

    const lines = hunk.patch.split('\n')

    for (const line of lines) {
      starts = ends
      ends += line.length + 1 // add trimmed newline

      // Conditions to keep forward

      if (!started) {
        if (line.startsWith('@@ ')) {
          started = true
        }
        continue
      }

      if (line.startsWith('@@ ')) {
        blockId++
        rowIndex = 0
        continue
      }

      if (line.includes('No newline at end of file')) {
        continue
      }

      const val: endsVal = {
        starts,
        ends,
        rowIndex,
        blockId,
      }

      if (!endsAt.has(hunk.hunkID)) {
        endsAt.set(hunk.hunkID, [val])
      } else {
        const ref = endsAt.get(hunk.hunkID)
        if (ref) {
          ref.push(val)
        }
      }

      rowIndex++
    }
  }

  for (const [hunkID, foundIndexes] of searchResult) {
    const ranges = endsAt.get(hunkID)

    if (!ranges) {
      continue
    }

    for (const foundIndex of foundIndexes) {
      for (let i = 0; i < ranges.length; i++) {
        const { starts, ends, rowIndex, blockId } = ranges[i]

        if (foundIndex < starts) {
          break
        }

        if (foundIndex < ends) {
          res.add(hunkID + '-' + rowIndex + '-' + blockId)
          break
        }
      }
    }
  }

  return res
}

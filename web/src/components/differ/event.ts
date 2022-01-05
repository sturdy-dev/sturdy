export interface DifferSetHunksWithPrefix {
  prefix: string | null
  selected: boolean
}

export interface Comment {
  id: string
  codebase_id: string
  change_id: string
  user_id: string
  author: Author
  message: string
  path: number
  line_start: number
  line_end: number
  line_is_new: boolean
  created_at: number
  deleted_at: number
}

export interface Author {
  id: string
  name: string
  avatar_url: string
  user_id: string
}

export interface Hunk {
  patch: string
  diff: Diff
}

export interface Block {
  header: string
  lines: Array<Line>
}

export interface Line {
  content: string
  type: string
  oldNumber: number
  newNumber: number
}

export interface HighlightedBlock {
  header: string
  lines: Array<HighlightedLine>
}

export interface HighlightedLine {
  prefix: string
  content: string
  originalContent: string
  type: string
  oldNumber: number
  newNumber: number
}

interface Diff {
  addedLines: number
  deletedLines: number
  blocks: Array<Block>
  checksumAfter: string
  checksumBefore: string
  isCombined: boolean
  isGitDiff: boolean

  isNew: boolean
  isDeleted: boolean
  isMoved: boolean

  language: string // File extension
  mode: string
  newName: string
  oldName: string
}

export type CommentState = {
  isExpanded: boolean
  composingReply: string | undefined
}

export type SetCommentExpandedEvent = {
  commentId: string
  isExpanded: boolean
}

export type SetCommentComposingReply = {
  commentId: string
  composingReply: string | undefined
}

export const temporaryNewCommentID = (
  path: string,
  lineStart: number,
  lineIsNew: boolean
): string => {
  return `${path}-${lineStart}-${lineIsNew}`
}

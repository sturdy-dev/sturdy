export interface SetFileIsHiddenEvent {
  fileKey: string
  isHidden: boolean
}

export type DifferState = {
  isHidden: boolean
}

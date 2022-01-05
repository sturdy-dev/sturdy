import { EmojiConvertor } from 'emoji-js'

let converter = new EmojiConvertor()
converter.colons_mode = true

export const ConvertEmojiToColons = function (str: string): string {
  return converter.replace_unified(str)
}

type User = {
  id: string
  name: string
}

// mentionify highlights the mentioned users in the text by id.
// txt: string - the text to highlight
// users: User[] - the users to highlight
// mention: string - the mention string to use
// style: string - the style to use for highlighting
// returns: string - the highlighted text
// example:
//   mentionify('@user1 @user2', [{id: 'user1', name: 'User 1'}, {id: 'user2', name: 'User 2'}], '@', 'mention')
// returns '<span class="mention">@User 1</span> <span class="mention">@user 2</span>'
const mentionify = function (txt: string, mention: string, users: User[], style?: string): string {
  // loop through each user
  for (let i = 0; i < users.length; i++) {
    // get the user
    const user = users[i]
    // get the mention string
    const mentionStr = mention + user.id
    const replaceWith =
      style !== undefined
        ? '$1<span class="' + style + '">' + mention + user.name + '</span>'
        : '$1' + mention + user.name
    // replace the user's name with a span with the style. do not replace if word doesn't start with the mention string
    txt = txt.replace(new RegExp(`(^|\\s)${mentionStr}($|\\b)`, 'g'), replaceWith)
  }
  return txt
}

export default mentionify

package decorate

import (
	"fmt"
	"regexp"
	"strings"

	"mash/pkg/user"
)

// mentionRegexp returns a regexp that matches a @mention.
func mentionRegexp(mention string) *regexp.Regexp {
	// \s is a space or newline.
	// \b is a word boundary.
	return regexp.MustCompile(fmt.Sprintf("(^|\\s)@%s($|\\b)", mention))
}

// ExtractEmailMentions returns a map of user id mentions to user in the comment message.
func ExtractIDMentions(msg string, users []*user.User) map[string]*user.User {
	mentions := make(map[string]*user.User)
	for _, user := range users {
		re := mentionRegexp(user.ID)
		match := re.FindString(msg)
		if match == "" {
			continue
		}
		match = strings.TrimSpace(match)
		mentions[match] = user
	}
	return mentions
}

// ExtractNameMentions returns a map of user name mentions to user in the comment message.
func ExtractNameMentions(msg string, users []*user.User) map[string]*user.User {
	mentions := make(map[string]*user.User)
	for _, user := range users {
		re := mentionRegexp(user.Name)
		match := re.FindString(msg)
		if match == "" {
			continue
		}
		match = strings.TrimSpace(match)
		mentions[match] = user
	}
	return mentions
}

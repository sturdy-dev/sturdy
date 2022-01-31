package decorate

import (
	"fmt"
	"regexp"
	"strings"

	"getsturdy.com/api/pkg/users"
)

// mentionRegexp returns a regexp that matches a @mention.
func mentionRegexp(mention string) *regexp.Regexp {
	// \s is a space or newline.
	// \b is a word boundary.
	return regexp.MustCompile(fmt.Sprintf("(^|\\s)@%s($|\\b)", mention))
}

// ExtractEmailMentions returns a map of user id mentions to user in the comment message.
func ExtractIDMentions(msg string, uu []*users.User) map[string]*users.User {
	mentions := make(map[string]*users.User)
	for _, user := range uu {
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
func ExtractNameMentions(msg string, uu []*users.User) map[string]*users.User {
	mentions := make(map[string]*users.User)
	for _, user := range uu {
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

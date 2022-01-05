package message

import (
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func CommitMessage(draftDescription string) string {
	newLiner := strings.NewReplacer(
		"<ul>", "<ul>\n",
		"<ol>", "<ol>\n",
		"</ol>", "</ol>\n",
		"<li><p>", "<li><p>\n* ",
		"<h1>", "\n<h1>\n",
		"<h2>", "\n<h2>\n",
		"<h3>", "\n<h3>\n",
		"<h4>", "\n<h4>\n",
		"<h5>", "\n<h5>\n",
		"<h6>", "\n<h6>\n",
		"<br>", "<br>\n",
		"<p>", "<p>\n",
	)

	sanitized := bluemonday.StrictPolicy().Sanitize(newLiner.Replace(draftDescription))

	return strings.TrimSpace(
		// bluemonday normalizes all newlines to \n
		// We want to have Windows-compatible newlines, so replace all \n with \r\n
		strings.ReplaceAll(sanitized, "\n", "\r\n"),
	)
}

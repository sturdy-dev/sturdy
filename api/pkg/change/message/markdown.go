package message

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
)

func MarkdownToHtml(source string) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(source), &buf); err != nil {
		return "", fmt.Errorf("failed to convert markdown to html: %w", err)
	}
	return buf.String(), nil
}

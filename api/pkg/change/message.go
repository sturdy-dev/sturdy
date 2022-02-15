package change

import (
	"sort"
	"strings"
)

// ChangeMetadata is change data written to the commit messages
type ChangeMetadata struct {
	ChangeID         string `json:"change_id,omitempty"` // Used to keep track of rebases etc.
	Description      string `json:"description,omitempty"`
	UserID           string `json:"user_id,omitempty"`
	ViewID           string `json:"view_id,omitempty"`
	RevertedChangeID string `json:"reverted_change_id,omitempty"` // If this is a revert, points to the change that was reverted.
	WorkspaceID      string `json:"workspace_id,omitempty"`
}

func (c ChangeMetadata) ToCommitMessage() string {
	lines := []string{
		c.Description,
		"", // extra empty line
		"--- sturdy ---",
	}

	var metaLines []string

	// add all non-empty values
	for k, v := range map[string]string{
		"change_id":          c.ChangeID,
		"user_id":            c.UserID,
		"view_id":            c.ViewID,
		"reverted_change_id": c.RevertedChangeID,
		"workspace_id":       c.WorkspaceID,
	} {
		if v == "" {
			continue
		}
		metaLines = append(metaLines, k+": "+v)
	}

	// Keep ordering consistent
	sort.Strings(metaLines)
	lines = append(lines, metaLines...)

	return strings.Join(lines, "\r\n")
}

func ParseCommitMessage(input string) ChangeMetadata {
	// If this was a squash merge made via GitHub, GitHub replaces unix new lines with a Windows new line \r\n
	// in := strings.Replace(input, "\r\n", "\n", -1)

	parts := strings.Split(input, "\r\n--- sturdy ---\r\n")

	var res ChangeMetadata
	res.Description = strings.TrimSpace(parts[0])

	// this commit has no sturdy metadata
	if len(parts) != 2 {
		return res
	}

	// parse the key-values
	metadata := strings.TrimSpace(parts[1])
	for _, row := range strings.Split(metadata, "\r\n") {
		kv := strings.Split(row, ": ")
		if len(kv) != 2 {
			continue
		}

		switch kv[0] {
		case "change_id":
			res.ChangeID = kv[1]
		case "user_id":
			res.UserID = kv[1]
		case "view_id":
			res.ViewID = kv[1]
		case "reverted_change_id":
			res.RevertedChangeID = kv[1]
		case "workspace_id":
			res.WorkspaceID = kv[1]
		}
	}

	return res
}

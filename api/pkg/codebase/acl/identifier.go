package acl

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/match"
)

type Identifier struct {
	Type    identityType `json:"type,omitempty"`
	Pattern string       `json:"pattern,omitempty"`
}

// MarshalJSON implements encoding/json.Marshaller to override resulting format.
func (i *Identifier) MarshalJSON() ([]byte, error) {
	if i.Pattern == "" && i.Type == "" {
		return json.Marshal("")
	}
	if i.Pattern == "" {
		return json.Marshal(i.Type)
	}
	if i.Type == "" {
		return json.Marshal(i.Pattern)
	}
	if i.Type == Users {
		return json.Marshal(i.Pattern)
	}
	return json.Marshal(fmt.Sprintf("%s::%s", i.Type, i.Pattern))
}

// UnmarshalJSON implements encoding/json.UnmarshalJSON to parse source JSON in a different way.
func (i *Identifier) UnmarshalJSON(v []byte) error {
	s := new(string)
	if err := json.Unmarshal(v, s); err != nil {
		return err
	}

	parts := strings.SplitN(*s, "::", 2)
	if i == nil {
		i = new(Identifier)
	}

	if len(parts) > 1 {
		i.Type = identityType(parts[0])
		i.Pattern = parts[1]
	} else if len(parts) == 1 {
		i.Type = Users
		i.Pattern = parts[0]
	}
	return nil
}

func (i *Identifier) Matches(identity Identity) bool {
	if i.Type != identity.Type {
		return false
	}
	return match.Match(identity.ID, i.Pattern)
}

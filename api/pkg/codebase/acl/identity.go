package acl

import (
	"encoding/json"
	"fmt"
	"strings"
)

type identityType string

var supportedIdentityTypes = map[identityType]bool{
	Users:     true,
	Groups:    true,
	Codebases: true,
	ACLs:      true,
	Files:     true,
}

func (it identityType) IsValid() bool {
	return supportedIdentityTypes[it]
}

const (
	Users     identityType = "users"
	Codebases identityType = "codebases"
	Groups    identityType = "groups"
	ACLs      identityType = "acls"
	Files     identityType = "files"
)

type Identity struct {
	ID   string       `json:"id,omitempty"`
	Type identityType `json:"type,omitempty"`
}

// MarshalJSON implements encoding/json.Marshaller to override resulting format.
func (i *Identity) MarshalJSON() ([]byte, error) {
	if i.ID == "" && i.Type == "" {
		return json.Marshal("")
	}
	if i.ID == "" {
		return json.Marshal(i.Type)
	}
	if i.Type == "" {
		return json.Marshal(i.ID)
	}
	if i.Type == Users {
		return json.Marshal(i.ID)
	}
	return json.Marshal(fmt.Sprintf("%s::%s", i.Type, i.ID))
}

func (i *Identity) ParseString(s string) {
	parts := strings.SplitN(s, "::", 2)
	if i == nil {
		i = new(Identity)
	}

	if len(parts) > 1 {
		i.ID = parts[1]
		i.Type = identityType(parts[0])
	} else if len(parts) == 1 {
		i.ID = parts[0]
		i.Type = Users
	}
}

// UnmarshalJSON implements encoding/json.UnmarshalJSON to parse source JSON in a different way.
func (i *Identity) UnmarshalJSON(v []byte) error {
	s := new(string)
	if err := json.Unmarshal(v, s); err != nil {
		return err
	}
	i.ParseString(*s)
	return nil
}

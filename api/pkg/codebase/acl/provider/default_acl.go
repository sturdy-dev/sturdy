package provider

import (
	"bytes"
	"text/template"
)

var (
	defaultACLTemplate = template.Must(template.New("").Parse(`{
  "rules": [
    {
      "id": "everyone can manage access control",
      "principals": ["groups::everyone"],
      "action": "write",
      "resources": ["acls::{{ .aclID }}"],
    },
    {
      "id": "everyone can access all files",
      "principals": ["groups::everyone"],
      "action": "write",
      "resources": ["files::*"],
    },
  ],
  "groups": [
    {
      "id": "everyone",
      "members": ["*"],
    },
  ],
  "tests": [
    {
      /*
        make sure that at least someone can manage access control lists
      */
      "id": "{{ index .userEmails 0 }} can manage access control",
      "principal": "{{ index .userEmails 0}}",
      "allow": "write",
      "resource": "acls::{{ .aclID }}",
    },
  ],
}`))
)

func defaultACLFor(aclID string, userEmails []string) (string, error) {
	w := bytes.Buffer{}

	if err := defaultACLTemplate.Execute(&w, map[string]any{
		"aclID":      aclID,
		"userEmails": userEmails,
	}); err != nil {
		return "", err
	}

	return w.String(), nil
}

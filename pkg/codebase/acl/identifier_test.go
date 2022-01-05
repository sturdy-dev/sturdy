package acl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Identifier_json(t *testing.T) {
	testCases := []struct {
		String     string
		Identifier Identifier
	}{
		{`"user@example.org"`, Identifier{Type: Users, Pattern: "user@example.org"}},
		{`"groups::admins"`, Identifier{Type: Groups, Pattern: "admins"}},
	}

	for _, test := range testCases {
		t.Run(test.String, func(t *testing.T) {
			marshaled, err := test.Identifier.MarshalJSON()
			assert.NoError(t, err, "failed to marshal identifier")
			assert.Equal(t, test.String, string(marshaled))

			unmarshaled := Identifier{}
			assert.NoError(t, unmarshaled.UnmarshalJSON([]byte(test.String)), "failed to unmarshal identifier")
			assert.Equal(t, test.Identifier, unmarshaled)
		})
	}
}

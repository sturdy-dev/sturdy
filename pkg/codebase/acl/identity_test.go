package acl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Identity_json(t *testing.T) {
	testCases := []struct {
		String   string
		Identity Identity
	}{
		{`"user@example.org"`, Identity{Type: Users, ID: "user@example.org"}},
		{`"groups::admins"`, Identity{Type: Groups, ID: "admins"}},
	}

	for _, test := range testCases {
		t.Run(test.String, func(t *testing.T) {
			marshaled, err := test.Identity.MarshalJSON()
			assert.NoError(t, err, "failed to marshal identifier")
			assert.Equal(t, test.String, string(marshaled))

			unmarshaled := Identity{}
			assert.NoError(t, unmarshaled.UnmarshalJSON([]byte(test.String)), "failed to unmarshal identifier")
			assert.Equal(t, test.Identity, unmarshaled)
		})
	}
}

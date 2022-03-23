package workspaces

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NameOrFallback(t *testing.T) {
	testCases := []struct {
		Name             *string
		DraftDescription string
		Out              string
	}{
		{Out: "Untitled draft"},
		{Name: pString("name"), Out: "name"},
		{DraftDescription: "Test", Out: "Test"},
		{Name: pString("name"), DraftDescription: "Test", Out: "Test"},
		{DraftDescription: "<p>Test</p>", Out: "Test"},
		{DraftDescription: "<p>Test</p><p>description</p>", Out: "Test"},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			assert.Equal(t, tc.Out, Workspace{
				Name:             tc.Name,
				DraftDescription: tc.DraftDescription,
			}.NameOrFallback())
		})
	}
}

func pString(s string) *string {
	return &s
}

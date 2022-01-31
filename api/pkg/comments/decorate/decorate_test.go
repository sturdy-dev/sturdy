package decorate_test

import (
	"testing"

	"getsturdy.com/api/pkg/comments/decorate"
	"getsturdy.com/api/pkg/users"

	"github.com/stretchr/testify/assert"
)

var usersByID = map[string]*users.User{
	"user1": {
		ID:   "user1",
		Name: "User 1",
	},
	"user2": {
		ID:   "user2",
		Name: "User 2",
	},
	"user3": {
		ID:   "user3",
		Name: "User 3",
	},
}

var uu = []*users.User{}

func init() {
	for _, user := range usersByID {
		uu = append(uu, user)
	}
}

func TestExtractNameMentions(t *testing.T) {
	testCases := []struct {
		str string
		exp map[string]*users.User
	}{
		{
			str: "@User 1",
			exp: map[string]*users.User{
				"@User 1": usersByID["user1"],
			},
		},
		{
			str: "Hello @User 1",
			exp: map[string]*users.User{
				"@User 1": usersByID["user1"],
			},
		},
		{
			str: "Hello @User 1!",
			exp: map[string]*users.User{
				"@User 1": usersByID["user1"],
			},
		},
		{
			str: "Hello User 1!",
			exp: map[string]*users.User{},
		},
		{
			str: "Hello@User 1",
			exp: map[string]*users.User{},
		},
		{
			str: "Hello@User 1!",
			exp: map[string]*users.User{},
		},
		{
			str: "@User 1hello",
			exp: map[string]*users.User{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.str, func(t *testing.T) {
			actual := decorate.ExtractNameMentions(testCase.str, uu)
			assert.Equal(t, testCase.exp, actual)
		})
	}
}

func TestExtractIDMentions(t *testing.T) {
	testCases := []struct {
		str string
		exp map[string]*users.User
	}{
		{
			str: "@user1",
			exp: map[string]*users.User{
				"@user1": usersByID["user1"],
			},
		},
		{
			str: "Hello @user1",
			exp: map[string]*users.User{
				"@user1": usersByID["user1"],
			},
		},
		{
			str: "Hello @user1!",
			exp: map[string]*users.User{
				"@user1": usersByID["user1"],
			},
		},
		{
			str: "Hello user1!",
			exp: map[string]*users.User{},
		},
		{
			str: "Hello@user1",
			exp: map[string]*users.User{},
		},
		{
			str: "Hello@user1!",
			exp: map[string]*users.User{},
		},
		{
			str: "@user1hello",
			exp: map[string]*users.User{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.str, func(t *testing.T) {
			actual := decorate.ExtractIDMentions(testCase.str, uu)
			assert.Equal(t, testCase.exp, actual)
		})
	}
}

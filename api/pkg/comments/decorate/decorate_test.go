package decorate_test

import (
	"testing"

	"getsturdy.com/api/pkg/comments/decorate"
	"getsturdy.com/api/pkg/user"

	"github.com/stretchr/testify/assert"
)

var usersByID = map[string]*user.User{
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

var users = []*user.User{}

func init() {
	for _, user := range usersByID {
		users = append(users, user)
	}
}

func TestExtractNameMentions(t *testing.T) {
	testCases := []struct {
		str string
		exp map[string]*user.User
	}{
		{
			str: "@User 1",
			exp: map[string]*user.User{
				"@User 1": usersByID["user1"],
			},
		},
		{
			str: "Hello @User 1",
			exp: map[string]*user.User{
				"@User 1": usersByID["user1"],
			},
		},
		{
			str: "Hello @User 1!",
			exp: map[string]*user.User{
				"@User 1": usersByID["user1"],
			},
		},
		{
			str: "Hello User 1!",
			exp: map[string]*user.User{},
		},
		{
			str: "Hello@User 1",
			exp: map[string]*user.User{},
		},
		{
			str: "Hello@User 1!",
			exp: map[string]*user.User{},
		},
		{
			str: "@User 1hello",
			exp: map[string]*user.User{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.str, func(t *testing.T) {
			actual := decorate.ExtractNameMentions(testCase.str, users)
			assert.Equal(t, testCase.exp, actual)
		})
	}
}

func TestExtractIDMentions(t *testing.T) {
	testCases := []struct {
		str string
		exp map[string]*user.User
	}{
		{
			str: "@user1",
			exp: map[string]*user.User{
				"@user1": usersByID["user1"],
			},
		},
		{
			str: "Hello @user1",
			exp: map[string]*user.User{
				"@user1": usersByID["user1"],
			},
		},
		{
			str: "Hello @user1!",
			exp: map[string]*user.User{
				"@user1": usersByID["user1"],
			},
		},
		{
			str: "Hello user1!",
			exp: map[string]*user.User{},
		},
		{
			str: "Hello@user1",
			exp: map[string]*user.User{},
		},
		{
			str: "Hello@user1!",
			exp: map[string]*user.User{},
		},
		{
			str: "@user1hello",
			exp: map[string]*user.User{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.str, func(t *testing.T) {
			actual := decorate.ExtractIDMentions(testCase.str, users)
			assert.Equal(t, testCase.exp, actual)
		})
	}
}

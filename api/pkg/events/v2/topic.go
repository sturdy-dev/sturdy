package events

import (
	"fmt"

	"getsturdy.com/api/pkg/users"
)

type Topic string

func (t Topic) String() string {
	return string(t)
}

func userTopic(userID users.ID) Topic {
	return Topic(fmt.Sprintf("user:%s", userID))
}

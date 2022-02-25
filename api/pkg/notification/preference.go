package notification

import "getsturdy.com/api/pkg/users"

type Channel string

const (
	ChannelUndefined Channel = ""
	ChannelWeb       Channel = "web"
	ChannelEmail     Channel = "email"
)

// Preference is used to determine if user with _UserID_ wants to receive notifications of type _Type_ via _Channel_.
type Preference struct {
	UserID  users.ID         `db:"user_id"`
	Type    NotificationType `db:"type"`
	Channel Channel          `db:"channel"`
	Enabled bool             `db:"enabled"`
}

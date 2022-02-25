package newsletter

import "getsturdy.com/api/pkg/users"

type NotificationSettings struct {
	UserID            users.ID `db:"user_id"`
	ReceiveNewsletter bool     `db:"receive_newsletter"`
}

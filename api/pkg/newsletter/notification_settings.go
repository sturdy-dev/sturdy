package newsletter

type NotificationSettings struct {
	UserID            string `db:"user_id"`
	ReceiveNewsletter bool   `db:"receive_newsletter"`
}

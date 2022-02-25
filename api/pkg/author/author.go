package author

import "getsturdy.com/api/pkg/users"

type Author struct {
	UserID         users.ID `json:"user_id"`
	Name           string   `json:"name"`
	AvatarURL      string   `json:"avatar_url"`
	Email          string   `json:"email"`
	IsExternalUser bool     `json:"is_external_user"` // If the user is imported from Git
}

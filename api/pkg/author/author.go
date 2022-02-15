package author

type Author struct {
	UserID         string `json:"user_id"`
	Name           string `json:"name"`
	AvatarURL      string `json:"avatar_url"`
	Email          string `json:"email"`
	IsExternalUser bool   `json:"is_external_user"` // If the user is imported from Git
}

package author

import (
	"fmt"
	db_user "mash/pkg/user/db"
)

type Author struct {
	UserID         string `json:"user_id"`
	Name           string `json:"name"`
	AvatarURL      string `json:"avatar_url"`
	Email          string `json:"email"`
	IsExternalUser bool   `json:"is_external_user"` // If the user is imported from Git
}

func GetAuthor(userID string, userRepo db_user.Repository) (Author, error) {
	user, err := userRepo.Get(userID)
	if err != nil {
		return Author{}, fmt.Errorf("failed to get user %s: %w", userID, err)
	}

	return Author{
		UserID:    user.ID,
		Name:      user.Name,
		Email:     user.Email,
		AvatarURL: emptyIfNull(user.AvatarURL),
	}, nil
}

func emptyIfNull(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

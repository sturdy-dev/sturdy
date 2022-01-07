package service

import (
	"context"
	"fmt"

	"mash/pkg/notification"
	db_notification "mash/pkg/notification/db"
)

type Preferences struct {
	preferencesRepo *db_notification.PreferenceRepository
}

func NewPreferences(
	preferencesRepo *db_notification.PreferenceRepository,
) *Preferences {
	return &Preferences{
		preferencesRepo: preferencesRepo,
	}
}

var (
	supportedChannels = map[notification.Channel]bool{
		notification.ChannelEmail: true,
		notification.ChannelWeb:   true,
	}
)

// Update updates or creates a preference.
func (s *Preferences) Update(
	ctx context.Context,
	userID string,
	typ notification.NotificationType,
	channel notification.Channel,
	enabled bool,
) (*notification.Preference, error) {
	p := &notification.Preference{
		UserID:  userID,
		Type:    typ,
		Channel: channel,
		Enabled: enabled,
	}
	if err := s.preferencesRepo.Upsert(ctx, p); err != nil {
		return nil, fmt.Errorf("failed to upsert: %w", err)
	}
	return p, nil
}

// ListByUserID returns all existing user preferences from the database + default preferences for all other
// possible permutation of notification types and channels.
func (s *Preferences) ListByUserID(ctx context.Context, userID string) ([]*notification.Preference, error) {
	pp, err := s.preferencesRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	existing := map[notification.Channel]map[notification.NotificationType]*notification.Preference{
		notification.ChannelEmail: {},
		notification.ChannelWeb:   {},
	}
	for _, p := range pp {
		existing[p.Channel][p.Type] = p
	}

	resultLength := len(supportedTypes) * len(supportedChannels)
	result := make([]*notification.Preference, 0, resultLength)
	for channel, supported := range supportedChannels {
		if !supported {
			continue
		}

		for typ, supported := range supportedTypes {
			if !supported {
				continue
			}

			if p, found := existing[channel][typ]; found {
				result = append(result, p)
			} else {
				result = append(result, &notification.Preference{
					Channel: channel,
					Type:    typ,
					UserID:  userID,
					Enabled: true,
				})
			}
		}
	}

	return result, nil
}

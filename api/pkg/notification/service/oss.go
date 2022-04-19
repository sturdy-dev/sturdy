//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package service

import "getsturdy.com/api/pkg/notification"

var (
	supportedTypes = map[notification.NotificationType]bool{
		notification.CommentNotificationType:         true,
		notification.ReviewNotificationType:          true,
		notification.RequestedReviewNotificationType: true,
		notification.NewSuggestionNotificationType:   true,
		notification.InvitedToCodebase:               true,
		notification.InvitedToOrganization:           true,
	}
	supportedChannels = map[notification.Channel]bool{
		notification.ChannelWeb: true,
	}
)

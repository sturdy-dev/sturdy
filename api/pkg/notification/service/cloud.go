//go:build cloud
// +build cloud

package service

import "getsturdy.com/api/pkg/notification"

var (
	supportedTypes = map[notification.NotificationType]bool{
		notification.CommentNotificationType:         true,
		notification.ReviewNotificationType:          true,
		notification.RequestedReviewNotificationType: true,
		notification.NewSuggestionNotificationType:   true,
		notification.GitHubRepositoryImported:        true,
	}
	supportedChannels = map[notification.Channel]bool{
		notification.ChannelEmail: true,
		notification.ChannelWeb:   true,
	}
)

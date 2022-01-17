//go:build enterprise || cloud
// +build enterprise cloud

package service

import "mash/pkg/notification"

var supportedTypes = map[notification.NotificationType]bool{
	notification.CommentNotificationType:         true,
	notification.ReviewNotificationType:          true,
	notification.RequestedReviewNotificationType: true,
	notification.NewSuggestionNotificationType:   true,
	notification.GitHubRepositoryImported:        true,
}

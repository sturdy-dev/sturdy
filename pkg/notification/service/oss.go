// +build !enterprise

package service

import "mash/pkg/notification"

var supportedTypes = map[notification.NotificationType]bool{
	notification.CommentNotificationType:         true,
	notification.ReviewNotificationType:          true,
	notification.RequestedReviewNotificationType: true,
	notification.NewSuggestionNotificationType:   true,
}

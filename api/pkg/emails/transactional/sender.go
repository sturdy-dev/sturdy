package transactional

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"getsturdy.com/api/pkg/analytics"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	service_change "getsturdy.com/api/pkg/changes/service"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	"getsturdy.com/api/pkg/comments"
	db_comments "getsturdy.com/api/pkg/comments/db"
	decorate_comments "getsturdy.com/api/pkg/comments/decorate"
	"getsturdy.com/api/pkg/emails"
	"getsturdy.com/api/pkg/emails/transactional/templates"
	"getsturdy.com/api/pkg/jwt"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	db_newsletter "getsturdy.com/api/pkg/newsletter/db"
	"getsturdy.com/api/pkg/notification"
	service_notification "getsturdy.com/api/pkg/notification/service"
	db_review "getsturdy.com/api/pkg/review/db"
	"getsturdy.com/api/pkg/suggestions"
	db_suggestion "getsturdy.com/api/pkg/suggestions/db"
	"getsturdy.com/api/pkg/users"
	db_users "getsturdy.com/api/pkg/users/db"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"

	"go.uber.org/zap"
)

type EmailSender interface {
	SendWelcome(context.Context, *users.User) error
	SendNotification(context.Context, *users.User, *notification.Notification) error
	SendConfirmEmail(context.Context, *users.User) error
	SendMagicLink(context.Context, *users.User, string) error
}

var ErrNotSupported = errors.New("notification type not supported")

type Sender struct {
	logger *zap.Logger
	sender emails.Sender

	userRepo                       db_users.Repository
	codebaseUserRepo               db_codebases.CodebaseUserRepository
	commentsRepo                   db_comments.Repository
	codebaseRepo                   db_codebases.CodebaseRepository
	workspaceRepo                  db_workspaces.Repository
	suggestionRepo                 db_suggestion.Repository
	reviewRepo                     db_review.ReviewRepository
	notificationSettingsRepository db_newsletter.NotificationSettingsRepository

	jwtService    *service_jwt.Service
	changeService *service_change.Service

	notificationPreferences *service_notification.Preferences
	analyticsService        *service_analytics.Service
}

func New(
	logger *zap.Logger,
	sender emails.Sender,

	userRepo db_users.Repository,
	codebaseUserRepo db_codebases.CodebaseUserRepository,
	commentsRepo db_comments.Repository,
	codebaseRepo db_codebases.CodebaseRepository,
	workspaceRepo db_workspaces.Repository,
	suggestionRepo db_suggestion.Repository,
	reviewRepo db_review.ReviewRepository,
	notificationSettingsRepository db_newsletter.NotificationSettingsRepository,

	jwtService *service_jwt.Service,
	changeService *service_change.Service,

	notificationPreferences *service_notification.Preferences,

	analyticsService *service_analytics.Service,
) *Sender {
	return &Sender{
		logger: logger,
		sender: sender,

		userRepo:                       userRepo,
		codebaseUserRepo:               codebaseUserRepo,
		commentsRepo:                   commentsRepo,
		codebaseRepo:                   codebaseRepo,
		workspaceRepo:                  workspaceRepo,
		suggestionRepo:                 suggestionRepo,
		reviewRepo:                     reviewRepo,
		notificationSettingsRepository: notificationSettingsRepository,

		jwtService:    jwtService,
		changeService: changeService,

		notificationPreferences: notificationPreferences,
		analyticsService:        analyticsService,
	}
}

func (e *Sender) SendMagicLink(ctx context.Context, user *users.User, code string) error {
	title := fmt.Sprintf("[Sturdy] Confirmation code: %s", code)
	return e.Send(ctx, user, title, templates.MagicLinkTemplate, &templates.MagicLinkTemplateData{
		User: user,
		Code: code,
	})
}

func (e *Sender) SendConfirmEmail(ctx context.Context, user *users.User) error {
	token, err := e.jwtService.IssueToken(ctx, user.ID.String(), time.Hour, jwt.TokenTypeVerifyEmail)
	if err != nil {
		return fmt.Errorf("failed to issue jwt token: %w", err)
	}

	title := "[Sturdy] Confirm your email"
	return e.Send(ctx, user, title, templates.VerifyEmailTemplate, &templates.VerifyEmailTemplateData{
		User:  user,
		Token: token,
	})
}

func (e *Sender) shouldSendNotification(ctx context.Context, usr *users.User, notificationType notification.NotificationType) (bool, error) {
	shouldSendEmail, err := shouldSendEmail(e.notificationSettingsRepository, usr)
	if err != nil {
		return false, err
	}

	if !shouldSendEmail {
		return false, nil
	}

	pp, err := e.notificationPreferences.ListByUserID(ctx, usr.ID)
	if err != nil {
		return false, err
	}

	for _, preference := range pp {
		if preference.Channel != notification.ChannelEmail {
			continue
		}
		if preference.Type != notificationType {
			continue
		}
		return preference.Enabled, nil
	}

	return false, fmt.Errorf("notification preference for %s not found", notificationType)
}

func (e *Sender) SendNotification(ctx context.Context, usr *users.User, notif *notification.Notification) error {
	shouldSendNotification, err := e.shouldSendNotification(ctx, usr, notif.NotificationType)
	if err != nil {
		return err
	}

	if !shouldSendNotification {
		return nil
	}

	switch notif.NotificationType {
	case notification.CommentNotificationType:
		if err := e.sendCommentNotification(ctx, usr, comments.ID(notif.ReferenceID)); err != nil {
			return fmt.Errorf("failed to send comment notification: %w", err)
		}
		return nil
	case notification.NewSuggestionNotificationType:
		if err := e.sendNewSuggestionNotification(ctx, usr, suggestions.ID(notif.ReferenceID)); err != nil {
			return fmt.Errorf("failed to send new suggestion notification: %w", err)
		}
		return nil
	case notification.RequestedReviewNotificationType:
		if err := e.sendRequestedReviewNotification(ctx, usr, notif.ReferenceID); err != nil {
			return fmt.Errorf("failed to send requested review notification: %w", err)
		}
		return nil
	case notification.ReviewNotificationType:
		if err := e.sendReviewNotification(ctx, usr, notif.ReferenceID); err != nil {
			return fmt.Errorf("failed to send review notification: %w", err)
		}
		return nil
	default:
		return ErrNotSupported
	}
}

func (e *Sender) sendReviewNotification(ctx context.Context, usr *users.User, reviewID string) error {
	r, err := e.reviewRepo.Get(ctx, reviewID)
	if err != nil {
		return fmt.Errorf("failed to find review: %w", err)
	}

	author, err := e.userRepo.Get(r.UserID)
	if err != nil {
		return fmt.Errorf("failed to find author: %w", err)
	}

	w, err := e.workspaceRepo.Get(r.WorkspaceID)
	if err != nil {
		return fmt.Errorf("failed to find workspace: %w", err)
	}

	c, err := e.codebaseRepo.Get(r.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to find codebase: %w", err)
	}

	title := fmt.Sprintf("[Sturdy] %s sent you a review", author.Name)
	data := &templates.NotificationReviewTemplateData{
		User: usr,

		Author:    author,
		Review:    r,
		Workspace: w,
		Codebase:  c,
	}
	return e.Send(ctx, usr, title, templates.NotificationReviewTemplate, data)
}

func (e *Sender) sendRequestedReviewNotification(ctx context.Context, usr *users.User, reviewID string) error {
	r, err := e.reviewRepo.Get(ctx, reviewID)
	if err != nil {
		return fmt.Errorf("failed to find review: %w", err)
	}

	requestedBy, err := e.userRepo.Get(*r.RequestedBy)
	if err != nil {
		return fmt.Errorf("failed to find author: %w", err)
	}

	w, err := e.workspaceRepo.Get(r.WorkspaceID)
	if err != nil {
		return fmt.Errorf("failed to find workspace: %w", err)
	}

	c, err := e.codebaseRepo.Get(r.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to find codebase: %w", err)
	}

	title := fmt.Sprintf("[Sturdy] %s asked for your feedback", requestedBy.Name)
	data := &templates.NotificationRequestedReviewTemplateData{
		User: usr,

		RequestedBy: requestedBy,
		Workspace:   w,
		Codebase:    c,
	}
	return e.Send(ctx, usr, title, templates.NotificationRequestedReviewTemplate, data)
}

func (e *Sender) sendNewSuggestionNotification(ctx context.Context, usr *users.User, suggestionID suggestions.ID) error {
	s, err := e.suggestionRepo.GetByID(ctx, suggestionID)
	if err != nil {
		return fmt.Errorf("failed to find suggestion: %w", err)
	}

	author, err := e.userRepo.Get(s.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	workspace, err := e.workspaceRepo.Get(s.ForWorkspaceID)
	if err != nil {
		return fmt.Errorf("failed to get workspace: %w", err)
	}

	codebase, err := e.codebaseRepo.Get(workspace.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to get codebase: %w", err)
	}

	title := fmt.Sprintf("[Sturdy] New suggestion on %s", workspace.NameOrFallback())
	data := &templates.NotificationNewSuggestionTemplateData{
		User:      usr,
		Author:    author,
		Workspace: workspace,
		Codebase:  codebase,
	}
	return e.Send(ctx, usr, title, templates.NotificationNewSuggestionTemplate, data)
}

func (e *Sender) getUsersByCodebaseID(ctx context.Context, codebaseID string) ([]*users.User, error) {
	codebaseUsers, err := e.codebaseUserRepo.GetByCodebase(codebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get codebase users: %w", err)
	}
	userIDs := make([]users.ID, 0, len(codebaseUsers))
	for _, codebaseUser := range codebaseUsers {
		userIDs = append(userIDs, codebaseUser.UserID)
	}

	users, err := e.userRepo.GetByIDs(ctx, userIDs...)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

func (e *Sender) sendCommentNotification(ctx context.Context, usr *users.User, commentID comments.ID) error {
	comment, err := e.commentsRepo.Get(commentID)
	if err != nil {
		return fmt.Errorf("failed to find comment: %w", err)
	}
	author, err := e.userRepo.Get(comment.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	codebaseUsers, err := e.getUsersByCodebaseID(ctx, comment.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	mentions := decorate_comments.ExtractIDMentions(comment.Message, codebaseUsers)
	// replace all mentions with names
	for mention, user := range mentions {
		comment.Message = strings.ReplaceAll(comment.Message, mention, fmt.Sprintf("@%s", user.Name))
	}

	codebase, err := e.codebaseRepo.Get(comment.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to get codebase: %w", err)
	}

	data := &templates.NotificationCommentTemplateData{
		User: usr,

		Comment:  &comment,
		Author:   author,
		Codebase: codebase,
	}

	switch {
	case comment.ParentComment != nil: // replied to
		parentComment, err := e.commentsRepo.Get(*comment.ParentComment)
		if err != nil {
			return fmt.Errorf("failed to get parent comment: %w", err)
		}

		parentAuthor, err := e.userRepo.Get(parentComment.UserID)
		if err != nil {
			return fmt.Errorf("failed to get parent author: %w", err)
		}

		data.Parent = &templates.NotificationCommentTemplateData{
			Comment:  &parentComment,
			Author:   parentAuthor,
			Codebase: codebase, // assumption: comments are always in the same codebase
		}

		switch {
		case parentComment.ChangeID != nil:
			change, err := e.changeService.GetChangeByID(ctx, *parentComment.ChangeID)
			if err != nil {
				return fmt.Errorf("failed to get parent change: %w", err)
			}
			data.Parent.Change = change
			title := fmt.Sprintf(
				"[Sturdy] %s repied to %s's comment on %s",
				author.Name, parentAuthor.Name, *change.Title)
			return e.Send(ctx, usr, title, templates.NotificationCommentTemplate, data)
		case parentComment.WorkspaceID != nil:
			workspace, err := e.workspaceRepo.Get(*parentComment.WorkspaceID)
			if err != nil {
				return fmt.Errorf("failed to get parent workspace: %w", err)
			}
			data.Parent.Workspace = workspace
			title := fmt.Sprintf("[Sturdy] %s replied to %s's comment on %s",
				author.Name, parentAuthor.Name, *workspace.Name)
			return e.Send(ctx, usr, title, templates.NotificationCommentTemplate, data)
		default:
			title := fmt.Sprintf("[Sturdy] %s replied to %s's comment",
				author.Name, parentAuthor.Name)
			return e.Send(ctx, usr, title, templates.NotificationCommentTemplate, data)
		}
	case comment.ChangeID != nil:
		change, err := e.changeService.GetChangeByID(ctx, *comment.ChangeID)
		if err != nil {
			return fmt.Errorf("failed to get change: %w", err)
		}
		data.Change = change
		title := fmt.Sprintf("[Sturdy] %s commented on %s", author.Name, *change.Title)
		return e.Send(ctx, usr, title, templates.NotificationCommentTemplate, data)
	case comment.WorkspaceID != nil:
		workspace, err := e.workspaceRepo.Get(*comment.WorkspaceID)
		if err != nil {
			return fmt.Errorf("failed to get comment workspace: %w", err)
		}
		data.Workspace = workspace
		title := fmt.Sprintf("[Sturdy] %s commented on %s", author.Name, workspace.NameOrFallback())
		return e.Send(ctx, usr, title, templates.NotificationCommentTemplate, data)
	default:
		title := fmt.Sprintf("[Sturdy] %s commented", author.Name)
		return e.Send(ctx, usr, title, templates.NotificationCommentTemplate, data)
	}
}

func (e *Sender) SendWelcome(ctx context.Context, u *users.User) error {
	return e.Send(
		ctx,
		u,
		"Welcome to Sturdy! 🐣",
		templates.WelcomeTemplate,
		&templates.WelcomeTemplateData{
			User: u,
		})
}

func (e *Sender) Send(
	ctx context.Context,
	u *users.User,
	subject string,
	template templates.Template,
	data any,
) error {
	content, err := templates.Render(template, data)
	if err != nil {
		return fmt.Errorf("failed to render email: %w", err)
	}

	e.logger.Info(
		"sending email",
		zap.Stringer("user_id", u.ID),
		zap.String("template", string(template)),
	)

	if err := e.sender.Send(ctx, &emails.Email{
		To:      u.Email,
		Subject: subject,
		Html:    content,
	}); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	e.analyticsService.Capture(ctx, "email_sent", analytics.DistinctID(u.ID.String()),
		analytics.Property("template", string(template)),
		analytics.Property("subject", subject),
	)

	return nil
}

func shouldSendEmail(notificationSettingsRepository db_newsletter.NotificationSettingsRepository, u *users.User) (bool, error) {
	if !u.EmailVerified {
		return false, nil
	}
	settings, err := notificationSettingsRepository.GetByUser(u.ID)
	switch {
	case err == nil:
		return settings.ReceiveNewsletter, nil
	case errors.Is(err, sql.ErrNoRows):
		return true, nil
	default:
		return false, fmt.Errorf("failed to get notification settings: %w", err)
	}
}

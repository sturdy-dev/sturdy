package transactional

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	db_change "mash/pkg/change/db"
	db_codebase "mash/pkg/codebase/db"
	"mash/pkg/comments"
	db_comments "mash/pkg/comments/db"
	decorate_comments "mash/pkg/comments/decorate"
	"mash/pkg/emails"
	db_github "mash/pkg/github/db"
	"mash/pkg/jwt"
	service_jwt "mash/pkg/jwt/service"
	db_newsletter "mash/pkg/newsletter/db"
	"mash/pkg/notification"
	service_notification "mash/pkg/notification/service"
	db_review "mash/pkg/review/db"
	"mash/pkg/suggestions"
	db_suggestion "mash/pkg/suggestions/db"
	"mash/pkg/user"
	db_users "mash/pkg/user/db"
	db_workspace "mash/pkg/workspace/db"

	posthog "github.com/posthog/posthog-go"
	"go.uber.org/zap"
)

type EmailSender interface {
	SendWelcome(context.Context, *user.User) error
	SendNotification(context.Context, *user.User, *notification.Notification) error
	SendConfirmEmail(context.Context, *user.User) error
	SendMagicLink(context.Context, *user.User, string) error
}

type emailSender struct {
	logger *zap.Logger
	sender emails.Sender

	userRepo                       db_users.Repository
	codebaseUserRepo               db_codebase.CodebaseUserRepository
	commentsRepo                   db_comments.Repository
	changeRepo                     db_change.Repository
	codebaseRepo                   db_codebase.CodebaseRepository
	workspaceRepo                  db_workspace.Repository
	suggestionRepo                 db_suggestion.Repository
	reviewRepo                     db_review.ReviewRepository
	notificationSettingsRepository db_newsletter.NotificationSettingsRepository
	githubRepositoryRepo           db_github.GitHubRepositoryRepo

	jwtService *service_jwt.Service

	notificationPreferences *service_notification.Preferences
	posthogClient           posthog.Client
}

func New(
	logger *zap.Logger,
	sender emails.Sender,

	userRepo db_users.Repository,
	codebaseUserRepo db_codebase.CodebaseUserRepository,
	commentsRepo db_comments.Repository,
	changeRepo db_change.Repository,
	codebaseRepo db_codebase.CodebaseRepository,
	workspaceRepo db_workspace.Repository,
	suggestionRepo db_suggestion.Repository,
	reviewRepo db_review.ReviewRepository,
	notificationSettingsRepository db_newsletter.NotificationSettingsRepository,
	githubRepositoryRepo db_github.GitHubRepositoryRepo,

	jwtService *service_jwt.Service,

	notificationPreferences *service_notification.Preferences,

	posthogClient posthog.Client,
) EmailSender {
	return &emailSender{
		logger: logger,
		sender: sender,

		userRepo:                       userRepo,
		codebaseUserRepo:               codebaseUserRepo,
		commentsRepo:                   commentsRepo,
		changeRepo:                     changeRepo,
		codebaseRepo:                   codebaseRepo,
		workspaceRepo:                  workspaceRepo,
		suggestionRepo:                 suggestionRepo,
		reviewRepo:                     reviewRepo,
		notificationSettingsRepository: notificationSettingsRepository,
		githubRepositoryRepo:           githubRepositoryRepo,

		jwtService: jwtService,

		notificationPreferences: notificationPreferences,
		posthogClient:           posthogClient,
	}
}

func (e *emailSender) SendMagicLink(ctx context.Context, user *user.User, code string) error {
	title := fmt.Sprintf("[Sturdy] Confirmation code: %s", code)
	return e.send(ctx, user, title, MagicLinkTemplate, &MagicLinkTemplateData{
		User: user,
		Code: code,
	})
}

func (e *emailSender) SendConfirmEmail(ctx context.Context, usr *user.User) error {
	token, err := e.jwtService.IssueToken(ctx, usr.ID, time.Hour, jwt.TokenTypeVerifyEmail)
	if err != nil {
		return fmt.Errorf("failed to issue jwt token: %w", err)
	}

	title := "[Sturdy] Confirm your email"
	return e.send(ctx, usr, title, VerifyEmailTemplate, &VerifyEmailTemplateData{
		User:  usr,
		Token: token,
	})
}

func (e *emailSender) shouldSendNotification(ctx context.Context, usr *user.User, notificationType notification.NotificationType) (bool, error) {
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

func (e *emailSender) SendNotification(ctx context.Context, usr *user.User, notif *notification.Notification) error {
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
	case notification.GitHubRepositoryImported:
		if err := e.sendGitHubRepositoryImportedNotification(ctx, usr, notif.ReferenceID); err != nil {
			return fmt.Errorf("failed to send github repository imported notification: %w", err)
		}
		return nil
	default:
		e.logger.Warn("email notification not supported", zap.String("type", string(notif.NotificationType)))
		return nil
	}
}

func (e *emailSender) sendReviewNotification(ctx context.Context, usr *user.User, reviewID string) error {
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
	data := &NotificationReviewTemplateData{
		User: usr,

		Author:    author,
		Review:    r,
		Workspace: w,
		Codebase:  c,
	}
	return e.send(ctx, usr, title, NotificationReviewTemplate, data)
}

func (e *emailSender) sendRequestedReviewNotification(ctx context.Context, usr *user.User, reviewID string) error {
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
	data := &NotificationRequestedReviewTemplateData{
		User: usr,

		RequestedBy: requestedBy,
		Workspace:   w,
		Codebase:    c,
	}
	return e.send(ctx, usr, title, NotificationRequestedReviewTemplate, data)
}

func (e *emailSender) sendNewSuggestionNotification(ctx context.Context, usr *user.User, suggestionID suggestions.ID) error {
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
	data := &NotificationNewSuggestionTemplateData{
		User:      usr,
		Author:    author,
		Workspace: workspace,
		Codebase:  codebase,
	}
	return e.send(ctx, usr, title, NotificationNewSuggestionTemplate, data)
}

func (e *emailSender) getUsersByCodebaseID(ctx context.Context, codebaseID string) ([]*user.User, error) {
	codebaseUsers, err := e.codebaseUserRepo.GetByCodebase(codebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get codebase users: %w", err)
	}
	userIDs := make([]string, 0, len(codebaseUsers))
	for _, codebaseUser := range codebaseUsers {
		userIDs = append(userIDs, codebaseUser.UserID)
	}

	users, err := e.userRepo.GetByIDs(ctx, userIDs...)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

func (e *emailSender) sendCommentNotification(ctx context.Context, usr *user.User, commentID comments.ID) error {
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

	data := &NotificationCommentTemplateData{
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

		data.Parent = &NotificationCommentTemplateData{
			Comment:  &parentComment,
			Author:   parentAuthor,
			Codebase: codebase, // assumption: comments are always in the same codebase
		}

		switch {
		case parentComment.ChangeID != nil:
			change, err := e.changeRepo.Get(*parentComment.ChangeID)
			if err != nil {
				return fmt.Errorf("failed to get parent change: %w", err)
			}
			data.Parent.Change = &change
			title := fmt.Sprintf(
				"[Strudy] %s repied to %s's comment on %s",
				author.Name, parentAuthor.Name, *change.Title)
			return e.send(ctx, usr, title, NotificationCommentTemplate, data)
		case parentComment.WorkspaceID != nil:
			workspace, err := e.workspaceRepo.Get(*parentComment.WorkspaceID)
			if err != nil {
				return fmt.Errorf("failed to get parent workspace: %w", err)
			}
			data.Parent.Workspace = workspace
			title := fmt.Sprintf("[Sturdy] %s replied to %s's comment on %s",
				author.Name, parentAuthor.Name, *workspace.Name)
			return e.send(ctx, usr, title, NotificationCommentTemplate, data)
		default:
			title := fmt.Sprintf("[Sturdy] %s replied to %s's comment",
				author.Name, parentAuthor.Name)
			return e.send(ctx, usr, title, NotificationCommentTemplate, data)
		}
	case comment.ChangeID != nil:
		change, err := e.changeRepo.Get(*comment.ChangeID)
		if err != nil {
			return fmt.Errorf("failed to get change: %w", err)
		}
		data.Change = &change
		title := fmt.Sprintf("[Sturdy] %s commented on %s", author.Name, *change.Title)
		return e.send(ctx, usr, title, NotificationCommentTemplate, data)
	case comment.WorkspaceID != nil:
		workspace, err := e.workspaceRepo.Get(*comment.WorkspaceID)
		if err != nil {
			return fmt.Errorf("failed to get comment workspace: %w", err)
		}
		data.Workspace = workspace
		title := fmt.Sprintf("[Sturdy] %s commented on %s", author.Name, workspace.NameOrFallback())
		return e.send(ctx, usr, title, NotificationCommentTemplate, data)
	default:
		title := fmt.Sprintf("[Sturdy] %s commented", author.Name)
		return e.send(ctx, usr, title, NotificationCommentTemplate, data)
	}
}

func (e *emailSender) sendGitHubRepositoryImportedNotification(ctx context.Context, u *user.User, gitHubRepoID string) error {
	repo, err := e.githubRepositoryRepo.GetByID(gitHubRepoID)
	if err != nil {
		return fmt.Errorf("failed to get github reposotory: %w", err)
	}
	c, err := e.codebaseRepo.Get(repo.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to get codebase: %w", err)
	}
	return e.send(
		ctx,
		u,
		fmt.Sprintf("[Sturdy] \"%s\" is now ready", c.Name),
		NotificationGitHubRepositoryImportedTemplate,
		&NotificationGitHubRepositoryImportedTemplateData{
			GitHubRepo: repo,
			Codebase:   c,
			User:       u,
		})
}

func (e *emailSender) SendWelcome(ctx context.Context, u *user.User) error {
	return e.send(
		ctx,
		u,
		"Welcome to Sturdy! 🐣",
		WelcomeTemplate,
		&WelcomeTemplateData{
			User: u,
		})
}

func (e *emailSender) send(
	ctx context.Context,
	u *user.User,
	subject string,
	template Template,
	data interface{},
) error {
	content, err := Render(template, data)
	if err != nil {
		return fmt.Errorf("failed to render email: %w", err)
	}

	e.logger.Info(
		"sending email",
		zap.String("user_id", u.ID),
		zap.String("template", string(template)),
	)

	if err := e.posthogClient.Enqueue(posthog.Capture{
		DistinctId: u.ID,
		Event:      "email_sent",
		Properties: map[string]interface{}{
			"user_id":  u.ID,
			"template": string(template),
			"subject":  subject,
		},
	}); err != nil {
		e.logger.Error("failed to enqueue posthog event", zap.Error(err))
	}

	if err := e.sender.Send(ctx, &emails.Email{
		To:      u.Email,
		Subject: subject,
		Html:    content,
	}); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func shouldSendEmail(notificationSettingsRepository db_newsletter.NotificationSettingsRepository, u *user.User) (bool, error) {
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

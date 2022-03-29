package templates

import (
	"bytes"
	"embed"
	"encoding/base64"
	"fmt"
	"html/template"

	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/comments"
	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/jwt"
	"getsturdy.com/api/pkg/review"
	"getsturdy.com/api/pkg/users"
	"getsturdy.com/api/pkg/workspaces"
)

var (
	//go:embed output/*.html output/notification/*.html
	fs        embed.FS
	templates = template.Must(template.New("").Funcs(map[string]any{
		"ternary": func(vt any, vf any, v bool) any {
			if v {
				return vt
			}
			return vf
		},
		"defaultString": func(dft string, in ...string) string {
			if len(in) == 0 || in[0] == "" {
				return dft
			}
			return in[0]
		},
		"base64Encode": func(in string) string {
			return base64.URLEncoding.EncodeToString([]byte(in))
		},
	}).ParseFS(fs, "output/*.html", "output/notification/*.html"))
)

type Template string

const (
	WelcomeTemplate                              Template = "welcome.template.html"
	NotificationGitHubRepositoryImportedTemplate Template = "github_repository_imported.template.html"
	NotificationCommentTemplate                  Template = "comment.template.html"
	NotificationNewSuggestionTemplate            Template = "new_suggestion.template.html"
	NotificationRequestedReviewTemplate          Template = "requested_review.template.html"
	NotificationReviewTemplate                   Template = "review.template.html"
	VerifyEmailTemplate                          Template = "verify_email.template.html"
	MagicLinkTemplate                            Template = "magic_link.template.html"
	InviteNewUserTemplate                        Template = "invite_new_user.template.html"
)

type WelcomeTemplateData struct {
	User *users.User
}

type NotificationGitHubRepositoryImportedTemplateData struct {
	GitHubRepo *github.Repository
	Codebase   *codebases.Codebase
	User       *users.User
}

type NotificationCommentTemplateData struct {
	User *users.User

	Comment   *comments.Comment
	Author    *users.User
	Codebase  *codebases.Codebase
	Workspace *workspaces.Workspace
	Change    *changes.Change

	Parent *NotificationCommentTemplateData
}

type NotificationNewSuggestionTemplateData struct {
	User *users.User

	Author    *users.User
	Workspace *workspaces.Workspace
	Codebase  *codebases.Codebase
}

type NotificationRequestedReviewTemplateData struct {
	User *users.User

	RequestedBy *users.User
	Workspace   *workspaces.Workspace
	Codebase    *codebases.Codebase
}

type NotificationReviewTemplateData struct {
	User *users.User

	Author    *users.User
	Review    *review.Review
	Workspace *workspaces.Workspace
	Codebase  *codebases.Codebase
}

type MagicLinkTemplateData struct {
	User *users.User
	Code string
}

type VerifyEmailTemplateData struct {
	User *users.User

	Token *jwt.Token
}

type InviteNewUserTemplateData struct {
	InvitingUser *users.User
	InvitedUser  *users.User
	Codebase     *codebases.Codebase
}

func Render(template Template, data any) (string, error) {
	rendered := &bytes.Buffer{}
	if err := templates.ExecuteTemplate(rendered, string(template), data); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}
	return rendered.String(), nil
}

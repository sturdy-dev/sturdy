package transactional

import (
	"bytes"
	"embed"
	"encoding/base64"
	"fmt"
	"html/template"

	"mash/pkg/change"
	"mash/pkg/codebase"
	"mash/pkg/comments"
	"mash/pkg/github"
	"mash/pkg/jwt"
	"mash/pkg/review"
	"mash/pkg/user"
	"mash/pkg/workspace"
)

var (
	//go:embed templates/output/*.html templates/output/notification/*.html
	fs        embed.FS
	templates = template.Must(template.New("").Funcs(map[string]interface{}{
		"ternary": func(vt interface{}, vf interface{}, v bool) interface{} {
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
	}).ParseFS(fs, "templates/output/*.html", "templates/output/notification/*.html"))
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
)

type WelcomeTemplateData struct {
	User *user.User
}

type NotificationGitHubRepositoryImportedTemplateData struct {
	GitHubRepo *github.GitHubRepository
	Codebase   *codebase.Codebase
	User       *user.User
}

type NotificationCommentTemplateData struct {
	User *user.User

	Comment   *comments.Comment
	Author    *user.User
	Codebase  *codebase.Codebase
	Workspace *workspace.Workspace
	Change    *change.Change

	Parent *NotificationCommentTemplateData
}

type NotificationNewSuggestionTemplateData struct {
	User *user.User

	Author    *user.User
	Workspace *workspace.Workspace
	Codebase  *codebase.Codebase
}

type NotificationRequestedReviewTemplateData struct {
	User *user.User

	RequestedBy *user.User
	Workspace   *workspace.Workspace
	Codebase    *codebase.Codebase
}

type NotificationReviewTemplateData struct {
	User *user.User

	Author    *user.User
	Review    *review.Review
	Workspace *workspace.Workspace
	Codebase  *codebase.Codebase
}

type MagicLinkTemplateData struct {
	User *user.User
	Code string
}

type VerifyEmailTemplateData struct {
	User *user.User

	Token *jwt.Token
}

func Render(template Template, data interface{}) (string, error) {
	rendered := &bytes.Buffer{}
	if err := templates.ExecuteTemplate(rendered, string(template), data); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}
	return rendered.String(), nil
}

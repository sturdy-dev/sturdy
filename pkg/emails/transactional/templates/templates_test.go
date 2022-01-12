package templates

import (
	"os"
	"testing"

	"mash/pkg/change"
	"mash/pkg/codebase"
	"mash/pkg/comments"
	"mash/pkg/github"
	"mash/pkg/jwt"
	"mash/pkg/review"
	"mash/pkg/user"
	"mash/pkg/workspace"

	"github.com/stretchr/testify/assert"
)

func TestRenderWelcome(t *testing.T) {
	output, err := Render(WelcomeTemplate, WelcomeTemplateData{
		User: &user.User{
			Email: "test@email.com",
		},
	})

	// uncomment make a snapshot
	// os.WriteFile("testdata/welcome.html", []byte(output), 0666)

	assert.NoError(t, err)
	assert.Equal(t, mustReadFile(t, "testdata/welcome.html"), output)
}

func TestRenderGitHubRepositoryImported(t *testing.T) {
	output, err := Render(NotificationGitHubRepositoryImportedTemplate, NotificationGitHubRepositoryImportedTemplateData{
		GitHubRepo: &github.GitHubRepository{
			Name: "codebase",
		},
		Codebase: &codebase.Codebase{
			Name:            "imported-codebase",
			ShortCodebaseID: "123456",
		},
		User: &user.User{
			Name:  "Nikita",
			Email: "test@email.com",
		},
	})

	// uncomment to make a snapshot
	// os.WriteFile("testdata/notification/github_repository_imported.html", []byte(output), 0666)

	assert.NoError(t, err)
	assert.Equal(t, mustReadFile(t, "testdata/notification/github_repository_imported.html"), output)
}

func TestRenderNotificationComment_commented(t *testing.T) {
	output, err := Render(NotificationCommentTemplate, NotificationCommentTemplateData{
		User: &user.User{
			Email: "test@email.com",
		},
		Comment: &comments.Comment{
			Message: "This is my comment message",
		},
		Author: &user.User{
			Name: "User One",
		},
		Codebase: &codebase.Codebase{
			ShortCodebaseID: "short-id",
			Name:            "codebase",
		},
	})

	// uncomment to make a snapshot
	os.WriteFile("testdata/notification/comment_commented.html", []byte(output), 0666)

	assert.NoError(t, err)
	assert.Equal(t, mustReadFile(t, "testdata/notification/comment_commented.html"), output)
}

func TestRenderNotificationComment_commented_on_change(t *testing.T) {
	output, err := Render(NotificationCommentTemplate, NotificationCommentTemplateData{
		User: &user.User{
			Email: "test@email.com",
		},
		Comment: &comments.Comment{
			Message: "This is my comment message",
		},
		Change: &change.Change{
			ID:    "change-id",
			Title: strPointer("make tests pass"),
		},
		Author: &user.User{
			Name: "User One",
		},
		Codebase: &codebase.Codebase{
			ShortCodebaseID: "short-id",
			Name:            "codebase",
		},
	})

	// uncomment to make a snapshot
	// os.WriteFile("testdata/notification/comment_commented_on_change.html", []byte(output), 0666)

	assert.NoError(t, err)
	assert.Equal(t, mustReadFile(t, "testdata/notification/comment_commented_on_change.html"), output)
}

func TestRenderNotificationComment_commented_on_workspace(t *testing.T) {
	output, err := Render(NotificationCommentTemplate, NotificationCommentTemplateData{
		User: &user.User{
			Email: "test@email.com",
		},
		Comment: &comments.Comment{
			Message: "This is my comment message",
		},
		Workspace: &workspace.Workspace{
			ID:   "workspace-id",
			Name: strPointer("workspace"),
		},
		Author: &user.User{
			Name: "User One",
		},
		Codebase: &codebase.Codebase{
			ShortCodebaseID: "short-id",
			Name:            "codebase",
		},
	})

	// uncomment to make a snapshot
	// os.WriteFile("testdata/notification/comment_commented_on_workspace.html", []byte(output), 0666)

	assert.NoError(t, err)
	assert.Equal(t, mustReadFile(t, "testdata/notification/comment_commented_on_workspace.html"), output)
}

func TestRenderNotificationComment_replied_your(t *testing.T) {
	usr := &user.User{
		ID:    "id",
		Name:  "User one",
		Email: "test@email.com",
	}

	output, err := Render(NotificationCommentTemplate, NotificationCommentTemplateData{
		User: usr,
		Comment: &comments.Comment{
			Message: "This is my comment message",
		},
		Author: &user.User{
			Name: "another user",
		},
		Codebase: &codebase.Codebase{
			ShortCodebaseID: "short-id",
			Name:            "codebase",
		},
		Parent: &NotificationCommentTemplateData{
			Author: usr,
		},
	})

	// uncomment to make a snapshot
	// os.WriteFile("testdata/notification/comment_replied_your.html", []byte(output), 0666)

	assert.NoError(t, err)
	assert.Equal(t, mustReadFile(t, "testdata/notification/comment_replied_your.html"), output)
}

func TestRenderNotificationComment_replied(t *testing.T) {
	output, err := Render(NotificationCommentTemplate, NotificationCommentTemplateData{
		User: &user.User{
			ID:    "0",
			Email: "test@email.com",
		},
		Comment: &comments.Comment{
			Message: "This is my comment message",
		},
		Author: &user.User{
			ID:   "1",
			Name: "User One",
		},
		Codebase: &codebase.Codebase{
			ShortCodebaseID: "short-id",
			Name:            "codebase",
		},
		Parent: &NotificationCommentTemplateData{
			Author: &user.User{
				ID:   "2",
				Name: "User two",
			},
		},
	})

	// uncomment to make a snapshot
	// os.WriteFile("testdata/notification/comment_replied.html", []byte(output), 0666)

	assert.NoError(t, err)
	assert.Equal(t, mustReadFile(t, "testdata/notification/comment_replied.html"), output)
}

func TestRenderNotificationComment_replied_on_workspace(t *testing.T) {
	output, err := Render(NotificationCommentTemplate, NotificationCommentTemplateData{
		User: &user.User{
			ID:    "0",
			Email: "test@email.com",
		},
		Comment: &comments.Comment{
			Message: "This is my comment message",
		},
		Author: &user.User{
			ID:   "1",
			Name: "User One",
		},
		Codebase: &codebase.Codebase{
			ShortCodebaseID: "short-id",
			Name:            "codebase",
		},
		Parent: &NotificationCommentTemplateData{
			Author: &user.User{
				ID:   "2",
				Name: "User two",
			},
			Workspace: &workspace.Workspace{
				ID:   "workspace-id",
				Name: strPointer("workspace"),
			},
		},
	})

	// uncomment to make a snapshot
	// os.WriteFile("testdata/notification/comment_replied_on_workspace.html", []byte(output), 0666)

	assert.NoError(t, err)
	assert.Equal(t, mustReadFile(t, "testdata/notification/comment_replied_on_workspace.html"), output)
}

func TestRenderNotificationComment_replied_your_on_workspace(t *testing.T) {
	usr := &user.User{
		ID:    "1",
		Name:  "user one",
		Email: "user@one.com",
	}
	output, err := Render(NotificationCommentTemplate, NotificationCommentTemplateData{
		User: usr,
		Comment: &comments.Comment{
			Message: "This is my comment message",
		},
		Author: &user.User{
			ID:   "2",
			Name: "User Two",
		},
		Codebase: &codebase.Codebase{
			ShortCodebaseID: "short-id",
			Name:            "codebase",
		},
		Parent: &NotificationCommentTemplateData{
			Author: usr,
			Workspace: &workspace.Workspace{
				ID:   "workspace-id",
				Name: strPointer("workspace"),
			},
		},
	})

	// uncomment to make a snapshot
	// os.WriteFile("testdata/notification/comment_replied_your_on_workspace.html", []byte(output), 0666)

	assert.NoError(t, err)
	assert.Equal(t, mustReadFile(t, "testdata/notification/comment_replied_your_on_workspace.html"), output)
}

func TestRenderNotificationComment_replied_on_change(t *testing.T) {
	output, err := Render(NotificationCommentTemplate, NotificationCommentTemplateData{
		User: &user.User{
			ID:    "0",
			Email: "test@email.com",
		},
		Comment: &comments.Comment{
			Message: "This is my comment message",
		},
		Author: &user.User{
			ID:   "1",
			Name: "User One",
		},
		Codebase: &codebase.Codebase{
			ShortCodebaseID: "short-id",
			Name:            "codebase",
		},
		Parent: &NotificationCommentTemplateData{
			Author: &user.User{
				ID:   "2",
				Name: "User two",
			},
			Change: &change.Change{
				ID:    "change-id",
				Title: strPointer("change"),
			},
		},
	})

	// uncomment to make a snapshot
	// os.WriteFile("testdata/notification/comment_replied_on_change.html", []byte(output), 0666)

	assert.NoError(t, err)
	assert.Equal(t, mustReadFile(t, "testdata/notification/comment_replied_on_change.html"), output)
}

func TestRenderNotificationComment_replied_your_on_change(t *testing.T) {
	usr := &user.User{
		Name:  "me",
		ID:    "0",
		Email: "me@test.com",
	}
	output, err := Render(NotificationCommentTemplate, NotificationCommentTemplateData{
		User: usr,
		Comment: &comments.Comment{
			Message: "This is my comment message",
		},
		Author: &user.User{
			ID:   "1",
			Name: "User One",
		},
		Codebase: &codebase.Codebase{
			ShortCodebaseID: "short-id",
			Name:            "codebase",
		},
		Parent: &NotificationCommentTemplateData{
			Author: usr,
			Change: &change.Change{
				ID:    "change-id",
				Title: strPointer("change"),
			},
		},
	})

	// uncomment to make a snapshot
	// os.WriteFile("testdata/notification/comment_replied_your_on_change.html", []byte(output), 0666)

	assert.NoError(t, err)
	assert.Equal(t, mustReadFile(t, "testdata/notification/comment_replied_your_on_change.html"), output)
}

func TestRenderNotificationNewSuggestion(t *testing.T) {
	usr := &user.User{
		Name:  "me",
		ID:    "0",
		Email: "me@test.com",
	}

	output, err := Render(NotificationNewSuggestionTemplate, NotificationNewSuggestionTemplateData{
		User: usr,
		Author: &user.User{
			ID:   "1",
			Name: "User One",
		},
		Codebase: &codebase.Codebase{
			ShortCodebaseID: "short-id",
			Name:            "codebase",
		},
		Workspace: &workspace.Workspace{
			ID:   "workspace-id",
			Name: strPointer("Workspace"),
		},
	})

	// uncomment to make a snapshot
	// os.WriteFile("testdata/notification/new_suggestion.html", []byte(output), 0666)

	assert.NoError(t, err)
	assert.Equal(t, mustReadFile(t, "testdata/notification/new_suggestion.html"), output)
}

func TestRenderNotificationRequestedReview(t *testing.T) {
	usr := &user.User{
		Name:  "me",
		ID:    "0",
		Email: "me@test.com",
	}

	output, err := Render(NotificationRequestedReviewTemplate, NotificationRequestedReviewTemplateData{
		User: usr,
		RequestedBy: &user.User{
			ID:   "1",
			Name: "User One",
		},
		Codebase: &codebase.Codebase{
			ShortCodebaseID: "short-id",
			Name:            "codebase",
		},
		Workspace: &workspace.Workspace{
			ID:   "workspace-id",
			Name: strPointer("Workspace"),
		},
	})

	// uncomment to make a snapshot
	// os.WriteFile("testdata/notification/requested_review.html", []byte(output), 0666)

	assert.NoError(t, err)
	assert.Equal(t, mustReadFile(t, "testdata/notification/requested_review.html"), output)
}

func TestRenderNotificationReview_approved(t *testing.T) {
	usr := &user.User{
		Name:  "me",
		ID:    "0",
		Email: "me@test.com",
	}

	output, err := Render(NotificationReviewTemplate, NotificationReviewTemplateData{
		User: usr,
		Author: &user.User{
			ID:   "1",
			Name: "User One",
		},
		Review: &review.Review{
			Grade: review.ReviewGradeApprove,
		},
		Codebase: &codebase.Codebase{
			ShortCodebaseID: "short-id",
			Name:            "codebase",
		},
		Workspace: &workspace.Workspace{
			ID:   "workspace-id",
			Name: strPointer("Workspace"),
		},
	})

	// uncomment to make a snapshot
	// os.WriteFile("testdata/notification/review_approved.html", []byte(output), 0666)

	assert.NoError(t, err)
	assert.Equal(t, mustReadFile(t, "testdata/notification/review_approved.html"), output)
}

func TestRenderNotificationReview_rejected(t *testing.T) {
	usr := &user.User{
		Name:  "me",
		ID:    "0",
		Email: "me@test.com",
	}

	output, err := Render(NotificationReviewTemplate, NotificationReviewTemplateData{
		User: usr,
		Author: &user.User{
			ID:   "1",
			Name: "User One",
		},
		Review: &review.Review{
			Grade: review.ReviewGradeReject,
		},
		Codebase: &codebase.Codebase{
			ShortCodebaseID: "short-id",
			Name:            "codebase",
		},
		Workspace: &workspace.Workspace{
			ID:   "workspace-id",
			Name: strPointer("Workspace"),
		},
	})

	// uncomment to make a snapshot
	// os.WriteFile("testdata/notification/review_rejected.html", []byte(output), 0666)

	assert.NoError(t, err)
	assert.Equal(t, mustReadFile(t, "testdata/notification/review_rejected.html"), output)
}

func TestRenderVerifyEmail(t *testing.T) {
	usr := &user.User{
		Name:  "me",
		ID:    "0",
		Email: "me@test.com",
	}

	output, err := Render(VerifyEmailTemplate, VerifyEmailTemplateData{
		User:  usr,
		Token: &jwt.Token{Token: "jwt-token"},
	})

	// uncomment to make a snapshot
	// os.WriteFile("testdata/verify_email.html", []byte(output), 0666)

	assert.NoError(t, err)
	assert.Equal(t, mustReadFile(t, "testdata/verify_email.html"), output)
}

func mustReadFile(t *testing.T, filename string) string {
	content, err := os.ReadFile(filename)
	assert.NoError(t, err)
	return string(content)
}

func strPointer(str string) *string {
	return &str
}

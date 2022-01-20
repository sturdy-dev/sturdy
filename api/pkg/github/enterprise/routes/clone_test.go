package routes

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"mash/pkg/analytics/disabled"
	"mash/pkg/github"
	"mash/pkg/github/config"
	"mash/pkg/github/enterprise/client"
	"mash/pkg/github/enterprise/routes/installation"
	"mash/pkg/github/enterprise/routes/internal/mock_client"
	"mash/pkg/github/enterprise/routes/internal/mock_sender"
	"mash/pkg/github/enterprise/service"
	workers_github "mash/pkg/github/enterprise/workers"
	"mash/pkg/internal/inmemory"
	"mash/pkg/notification"
	events "mash/pkg/view/events"
	"mash/vcs/testutil/executor"
	"mash/vcs/testutil/history"

	"github.com/golang/mock/gomock"
	gh "github.com/google/go-github/v39/github"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

//go:generate mockgen -destination internal/mock_client/repositories_client_mock.go mash/pkg/github/client RepositoriesClient
//go:generate mockgen -destination internal/mock_sender/notification_sender_mock.go mash/pkg/notification/sender NotificationSender

func TestCloneSendsNotifications(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	gitHubRepositoryRepo := inmemory.NewInMemoryGitHubRepositoryRepo()
	gitHubInstallationRepo := inmemory.NewInMemoryGitHubInstallationRepository()
	gitHubUserRepo := inmemory.NewInMemoryGitHubUserRepo()
	codebaseRepo := inmemory.NewInMemoryCodebaseRepo()
	codebaseUserRepo := inmemory.NewInMemoryCodebaseUserRepo()
	workspaceRepo := inmemory.NewInMemoryWorkspaceRepo()
	postHogClient := disabled.NewClient()
	executorProvider := executor.TestingExecutorProvider(t)
	eventsSender := events.NewSender(codebaseUserRepo, workspaceRepo, events.NewInMemory())

	ctx := context.Background()

	ghInstallation := github.GitHubInstallation{
		InstallationID:         rand.Int63n(1_000_000),
		ID:                     uuid.NewString(),
		Owner:                  uuid.NewString(),
		CreatedAt:              time.Now(),
		HasWorkflowsPermission: true,
	}
	assert.NoError(t, gitHubInstallationRepo.Create(ghInstallation))

	ghUser := github.GitHubUser{Username: "TEST-SENDER-" + uuid.NewString(), UserID: uuid.NewString(), CreatedAt: time.Now().Add(time.Hour * -3)}
	assert.NoError(t, gitHubUserRepo.Create(ghUser))

	ghUserCollaborator1 := github.GitHubUser{Username: "TEST-COLLABORATOR1-" + uuid.NewString(), UserID: uuid.NewString(), CreatedAt: time.Now().Add(time.Hour * -3)}
	assert.NoError(t, gitHubUserRepo.Create(ghUserCollaborator1))

	// This user is new, and should not receive any notifications
	ghUserCollaborator2New := github.GitHubUser{Username: "TEST-COLLABORATOR2-NEW-" + uuid.NewString(), UserID: uuid.NewString(), CreatedAt: time.Now()}
	assert.NoError(t, gitHubUserRepo.Create(ghUserCollaborator2New))

	fakeGitHubRepoPath := history.CreateRepoWithRootCommit(t, executorProvider)

	ghRepo := &gh.Repository{
		ID:       i(rand.Int63n(1_000_000)),
		Name:     str(uuid.NewString()),
		CloneURL: str("file://" + fakeGitHubRepoPath),
	}

	event := &gh.InstallationEvent{
		Action: str("created"),
		Sender: &gh.User{
			Login: str(ghUser.Username),
		},
		Installation: &gh.Installation{
			ID: &ghInstallation.InstallationID,
		},
		Repositories: []*gh.Repository{ghRepo},
	}

	ctrl := gomock.NewController(t)

	// users and the app have different permissions, use different mock clients to mock the expected behaviours values for each

	personalGitHubRepositoriesClient := mock_client.NewMockRepositoriesClient(ctrl)

	collaborators := []*gh.User{
		{Login: str(ghUser.Username)},
		{Login: str(ghUserCollaborator1.Username)},
		{Login: str(ghUserCollaborator2New.Username)},
		{Login: str(ghRepo.GetName() + "-collaborator-3")},
		{Login: str(ghRepo.GetName() + "-collaborator-4")},
	}
	personalGitHubRepositoriesClient.EXPECT().ListCollaborators(gomock.Any(), ghInstallation.Owner, ghRepo.GetName(), gomock.Any()).Return(collaborators, &gh.Response{NextPage: 0}, nil)

	appGitHubRepositoriesClient := mock_client.NewMockRepositoriesClient(ctrl)
	appGitHubRepositoriesClient.EXPECT().GetByID(gomock.Any(), ghRepo.GetID()).Return(ghRepo, nil, nil).MinTimes(1)

	notificationSender := mock_sender.NewMockNotificationSender(ctrl)

	// Expect both users to get notifications sent to them
	notificationSender.EXPECT().User(gomock.Any() /*context*/, ghUser.UserID, gomock.Any() /*codebaseID*/, notification.GitHubRepositoryImported, gomock.Any() /*referenceID*/).Times(1)
	notificationSender.EXPECT().User(gomock.Any() /*context*/, ghUserCollaborator1.UserID, gomock.Any() /*codebaseID*/, notification.GitHubRepositoryImported, gomock.Any() /*referenceID*/).Times(1)
	// This user is new, and should not receive any notifications
	notificationSender.EXPECT().User(gomock.Any() /*context*/, ghUserCollaborator2New.UserID, gomock.Any() /*codebaseID*/, notification.GitHubRepositoryImported, gomock.Any() /*referenceID*/).Times(0)

	importer := service.ImporterQueue(workers_github.NopImporter())
	svc := new(service.Service)
	cloner := service.ClonerQueue(&synchronousCloner{svc})
	*svc = *service.New(
		logger,
		gitHubRepositoryRepo,
		gitHubInstallationRepo,
		gitHubUserRepo,
		nil,
		config.GitHubAppConfig{},
		clientProvider(appGitHubRepositoriesClient),
		personalClientProvider(personalGitHubRepositoriesClient),
		&importer,
		&cloner,
		nil,
		nil,
		codebaseUserRepo,
		codebaseRepo,
		executorProvider,
		nil,
		postHogClient,
		notificationSender,
		eventsSender,

		nil,
	)

	err := installation.HandleInstallationEvent(
		ctx,
		logger,
		event,
		gitHubInstallationRepo,
		gitHubRepositoryRepo,
		postHogClient,
		codebaseRepo,
		svc,
	)
	assert.NoError(t, err)
}

func clientProvider(repoClient client.RepositoriesClient) func(gitHubAppConfig config.GitHubAppConfig, installationID int64) (tokenClient *client.GitHubClients, jwtClient *client.GitHubClients, err error) {
	return func(gitHubAppConfig config.GitHubAppConfig, installationID int64) (tokenClient *client.GitHubClients, jwtClient *client.GitHubClients, err error) {
		return &client.GitHubClients{
				Repositories: repoClient,
				PullRequests: &fakeGitHubPullRequestClient{},
			},
			&client.GitHubClients{
				Apps: &fakeGitHubAppsClient{},
			}, nil
	}
}

func personalClientProvider(repoClient client.RepositoriesClient) func(token string) (*client.GitHubClients, error) {
	return func(token string) (*client.GitHubClients, error) {
		return &client.GitHubClients{
			Repositories: repoClient,
			PullRequests: &fakeGitHubPullRequestClient{},
		}, nil
	}
}

type fakeGitHubPullRequestClient struct {
	prs []*gh.PullRequest
}

func (f *fakeGitHubPullRequestClient) List(ctx context.Context, owner string, repo string, opts *gh.PullRequestListOptions) ([]*gh.PullRequest, *gh.Response, error) {
	panic("implement me (list prs)")
}

func (f *fakeGitHubPullRequestClient) Create(ctx context.Context, owner string, repo string, pull *gh.NewPullRequest) (*gh.PullRequest, *gh.Response, error) {
	panic("implement me (create prs)")
}

func (f *fakeGitHubPullRequestClient) Get(ctx context.Context, owner string, repo string, number int) (*gh.PullRequest, *gh.Response, error) {
	panic("implement me (get pr)")
}

func (f *fakeGitHubPullRequestClient) Edit(ctx context.Context, owner string, repo string, number int, pull *gh.PullRequest) (*gh.PullRequest, *gh.Response, error) {
	panic("implement me (edit pr)")
}

type fakeGitHubAppsClient struct{}

func (f *fakeGitHubAppsClient) CreateInstallationToken(ctx context.Context, id int64, opts *gh.InstallationTokenOptions) (*gh.InstallationToken, *gh.Response, error) {
	return &gh.InstallationToken{
		Token:        str("testingtoken"),
		ExpiresAt:    t(time.Now().Add(time.Hour * 3)),
		Permissions:  opts.Permissions,
		Repositories: nil,
	}, nil, nil
}

type fakeGitHubRepositoriesClient struct{}

func (f *fakeGitHubRepositoriesClient) Get(ctx context.Context, owner, repo string) (*gh.Repository, *gh.Response, error) {
	panic("implement me")
}
func (f *fakeGitHubRepositoriesClient) GetByID(ctx context.Context, id int64) (*gh.Repository, *gh.Response, error) {
	panic("implement me")
}
func (f *fakeGitHubRepositoriesClient) ListCollaborators(ctx context.Context, owner, repo string, opts *gh.ListCollaboratorsOptions) ([]*gh.User, *gh.Response, error) {
	panic("implement me")
}

func str(s string) *string {
	return &s
}

func t(in time.Time) *time.Time {
	return &in
}

func i(i int64) *int64 {
	return &i
}

type synchronousCloner struct {
	gitHubService *service.Service
}

func (s *synchronousCloner) Enqueue(ctx context.Context, event *github.CloneRepositoryEvent) error {
	return s.gitHubService.Clone(event.CodebaseID, event.InstallationID, event.GitHubRepositoryID, event.SenderUserID)
}

package webhooks

import (
	"encoding/json"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/github/api"
)

// this is copied from: https://github.com/google/go-github/blob/41cbf81925093b28eed72d7ef81577a7c6fdfd1b/github/event_types.go
// if we need more fields, they should be copied manually to have control over the structure size

var eventTypeMapping = map[string]string{
	"installation":              "InstallationEvent",
	"installation_repositories": "InstallationRepositoriesEvent",
	"pull_request":              "PullRequestEvent",
	"push":                      "PushEvent",
	"status":                    "StatusEvent",
	"workflow_job":              "WorkflowJobEvent",
}

func ParseWebHook(messageType string, payload []byte) (any, error) {
	eventType, ok := eventTypeMapping[messageType]
	if !ok {
		return nil, fmt.Errorf("unknown X-Github-Event in message: %v", messageType)
	}

	event := Event{
		Type:       &eventType,
		RawPayload: (*json.RawMessage)(&payload),
	}
	return event.ParsePayload()
}

type Event struct {
	Type       *string          `json:"type,omitempty"`
	Public     *bool            `json:"public,omitempty"`
	RawPayload *json.RawMessage `json:"payload,omitempty"`
	Repo       *api.Repository  `json:"repo,omitempty"`
	Actor      *api.User        `json:"actor,omitempty"`
	CreatedAt  *time.Time       `json:"created_at,omitempty"`
	ID         *string          `json:"id,omitempty"`
}

// ParsePayload parses the event payload. For recognized event types,
// a value of the corresponding struct type will be returned.
func (e *Event) ParsePayload() (payload any, err error) {
	switch *e.Type {
	case "InstallationEvent":
		payload = &InstallationEvent{}
	case "InstallationRepositoriesEvent":
		payload = &InstallationRepositoriesEvent{}
	case "PullRequestEvent":
		payload = &PullRequestEvent{}
	case "PushEvent":
		payload = &PushEvent{}
	case "StatusEvent":
		payload = &StatusEvent{}
	case "WorkflowJobEvent":
		payload = &WorkflowJobEvent{}
	}
	err = json.Unmarshal(*e.RawPayload, &payload)
	return payload, err
}

type PullRequestEvent struct {
	PullRequest  *api.PullRequest  `json:"pull_request,omitempty"`
	Repo         *api.Repository   `json:"repository,omitempty"`
	Installation *api.Installation `json:"installation,omitempty"`
	Sender       *api.User         `json:"sender,omitempty"`
}

func (pre *PullRequestEvent) GetSender() *api.User {
	if pre == nil {
		return nil
	}
	return pre.Sender
}

func (pre *PullRequestEvent) GetPullRequest() *api.PullRequest {
	if pre == nil {
		return nil
	}
	return pre.PullRequest
}

func (pre *PullRequestEvent) GetRepo() *api.Repository {
	if pre == nil {
		return nil
	}
	return pre.Repo
}

func (pre *PullRequestEvent) GetInstallation() *api.Installation {
	if pre == nil {
		return nil
	}
	return pre.Installation
}

type PushEvent struct {
	Ref *string `json:"ref,omitempty"`

	Repo         *PushEventRepository `json:"repository,omitempty"`
	Installation *api.Installation    `json:"installation,omitempty"`
}

func (pr *PushEvent) GetInstallation() *api.Installation {
	if pr == nil {
		return nil
	}
	return pr.Installation
}

func (pr *PushEvent) GetRepo() *PushEventRepository {
	if pr == nil {
		return nil
	}
	return pr.Repo
}

func (pr *PushEvent) GetRef() string {
	if pr == nil || pr.Ref == nil {
		return ""
	}
	return *pr.Ref
}

type PushEventRepository struct {
	ID       *int64  `json:"id,omitempty"`
	FullName *string `json:"full_name,omitempty"`
}

func (per *PushEventRepository) GetID() int64 {
	if per == nil || per.ID == nil {
		return 0
	}
	return *per.ID
}

func (per *PushEventRepository) GetFullName() string {
	if per == nil || per.FullName == nil {
		return ""
	}
	return *per.FullName
}

type InstallationEvent struct {
	// The action that was performed. Can be either "created", "deleted", "suspend", "unsuspend" or "new_permissions_accepted".
	Action       *string           `json:"action,omitempty"`
	Repositories []*api.Repository `json:"repositories,omitempty"`
	Sender       *api.User         `json:"sender,omitempty"`
	Installation *api.Installation `json:"installation,omitempty"`
	// TODO key "requester" is not covered
}

func (ie *InstallationEvent) GetInstallation() *api.Installation {
	if ie == nil {
		return nil
	}
	return ie.Installation
}

func (ie *InstallationEvent) GetAction() string {
	if ie == nil || ie.Action == nil {
		return ""
	}
	return *ie.Action
}

type InstallationRepositoriesEvent struct {
	RepositoriesAdded   []*api.Repository `json:"repositories_added,omitempty"`
	RepositoriesRemoved []*api.Repository `json:"repositories_removed,omitempty"`
	RepositorySelection *string           `json:"repository_selection,omitempty"`
	Installation        *api.Installation `json:"installation,omitempty"`
}

func (ire *InstallationRepositoriesEvent) GetInstallation() *api.Installation {
	if ire == nil {
		return nil
	}
	return ire.Installation
}

func (ire *InstallationRepositoriesEvent) GetRepositorySelection() string {
	if ire == nil || ire.RepositorySelection == nil {
		return ""
	}
	return *ire.RepositorySelection
}

type StatusEvent struct {
	Context     *string `json:"context,omitempty"`
	State       *string `json:"state,omitempty"`
	Description *string `json:"description,omitempty"`
	TargetURL   *string `json:"target_url,omitempty"`

	Repo         *api.Repository       `json:"repository,omitempty"`
	Installation *api.Installation     `json:"installation,omitempty"`
	Commit       *api.RepositoryCommit `json:"commit,omitempty"`
	UpdatedAt    *api.Timestamp        `json:"updated_at,omitempty"`
	CreatedAt    *api.Timestamp        `json:"created_at,omitempty"`
}

func (se *StatusEvent) GetContext() string {
	if se == nil || se.Context == nil {
		return ""
	}
	return *se.Context
}

func (se *StatusEvent) GetCreatedAt() api.Timestamp {
	if se == nil || se.CreatedAt == nil {
		return api.Timestamp{}
	}
	return *se.CreatedAt
}

func (se *StatusEvent) GetUpdatedAt() api.Timestamp {
	if se == nil || se.UpdatedAt == nil {
		return api.Timestamp{}
	}
	return *se.UpdatedAt
}

func (se *StatusEvent) GetState() string {
	if se == nil || se.State == nil {
		return ""
	}
	return *se.State
}

func (se *StatusEvent) GetCommit() *api.RepositoryCommit {
	if se == nil {
		return nil
	}
	return se.Commit
}

func (se *StatusEvent) GetInstallation() *api.Installation {
	if se == nil {
		return nil
	}
	return se.Installation
}

func (se *StatusEvent) GetRepo() *api.Repository {
	if se == nil {
		return nil
	}
	return se.Repo
}

type WorkflowJobEvent struct {
	WorkflowJob  *api.WorkflowJob  `json:"workflow_job,omitempty"`
	Repo         *api.Repository   `json:"repository,omitempty"`
	Installation *api.Installation `json:"installation,omitempty"`
}

func (wje *WorkflowJobEvent) GetInstallation() *api.Installation {
	if wje == nil {
		return nil
	}
	return wje.Installation
}

func (wje *WorkflowJobEvent) GetRepo() *api.Repository {
	if wje == nil {
		return nil
	}
	return wje.Repo
}

func (wje *WorkflowJobEvent) GetWorkflowJob() *api.WorkflowJob {
	if wje == nil {
		return nil
	}
	return wje.WorkflowJob
}

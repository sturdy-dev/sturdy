package webhooks

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// this is copied from: https://github.com/google/go-github/blob/41cbf81925093b28eed72d7ef81577a7c6fdfd1b/github/event_types.go
// if we need more fields, they should be copied manually to have controll over the structure size

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
	Repo       *Repository      `json:"repo,omitempty"`
	Actor      *User            `json:"actor,omitempty"`
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
	PullRequest  *PullRequest  `json:"pull_request,omitempty"`
	Repo         *Repository   `json:"repository,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

func (pre *PullRequestEvent) GetPullRequest() *PullRequest {
	if pre == nil {
		return nil
	}
	return pre.PullRequest
}

func (pre *PullRequestEvent) GetRepo() *Repository {
	if pre == nil {
		return nil
	}
	return pre.Repo
}

func (pre *PullRequestEvent) GetInstallation() *Installation {
	if pre == nil {
		return nil
	}
	return pre.Installation
}

type PullRequest struct {
	ID             *int64     `json:"id,omitempty"`
	State          *string    `json:"state,omitempty"`
	Merged         *bool      `json:"merged,omitempty"`
	MergeCommitSHA *string    `json:"merge_commit_sha,omitempty"`
	HTMLURL        *string    `json:"html_url,omitempty"`
	ClosedAt       *time.Time `json:"closed_at,omitempty"`
	MergedAt       *time.Time `json:"merged_at,omitempty"`

	Base *PullRequestBranch `json:"base,omitempty"`
}

func (pr *PullRequest) GetID() int64 {
	if pr == nil || pr.ID == nil {
		return 0
	}
	return *pr.ID
}

func (pr *PullRequest) GetHTMLURL() string {
	if pr == nil || pr.HTMLURL == nil {
		return ""
	}
	return *pr.HTMLURL
}

func (pr *PullRequest) GetState() string {
	if pr == nil || pr.State == nil {
		return ""
	}
	return *pr.State
}

func (pr *PullRequest) GetMerged() bool {
	if pr == nil || pr.Merged == nil {
		return false
	}
	return *pr.Merged
}

func (pr *PullRequest) GetMergeCommitSHA() string {
	if pr == nil || pr.MergeCommitSHA == nil {
		return ""
	}
	return *pr.MergeCommitSHA
}

func (pr *PullRequest) GetBase() *PullRequestBranch {
	if pr == nil {
		return nil
	}
	return pr.Base
}

func (pr *PullRequest) GetClosedAt() *time.Time {
	if pr == nil || pr.ClosedAt == nil {
		return nil
	}
	return pr.ClosedAt
}

func (pr *PullRequest) GetMergedAt() *time.Time {
	if pr == nil || pr.MergedAt == nil {
		return nil
	}
	return pr.MergedAt
}

type PullRequestBranch struct {
	Label *string     `json:"label,omitempty"`
	Ref   *string     `json:"ref,omitempty"`
	SHA   *string     `json:"sha,omitempty"`
	Repo  *Repository `json:"repo,omitempty"`
	User  *User       `json:"user,omitempty"`
}

func (prb *PullRequestBranch) GetSHA() string {
	if prb == nil || prb.SHA == nil {
		return ""
	}
	return *prb.SHA
}

type Repository struct {
	ID       *int64  `json:"id,omitempty"`
	Name     *string `json:"name,omitempty"`
	FullName *string `json:"full_name,omitempty"`
}

func (r *Repository) GetID() int64 {
	if r == nil || r.ID == nil {
		return 0
	}
	return *r.ID
}

func (r *Repository) GetName() string {
	if r == nil || r.Name == nil {
		return ""
	}
	return *r.Name
}

func (r *Repository) GetFullName() string {
	if r == nil || r.FullName == nil {
		return ""
	}
	return *r.FullName
}

type Installation struct {
	ID      *int64 `json:"id,omitempty"`
	Account *User  `json:"account,omitempty"`
}

func (i *Installation) GetID() int64 {
	if i == nil || i.ID == nil {
		return 0
	}
	return *i.ID
}

func (i *Installation) GetAccount() *User {
	if i == nil {
		return nil
	}
	return i.Account
}

type PushEvent struct {
	Ref *string `json:"ref,omitempty"`

	Repo         *PushEventRepository `json:"repository,omitempty"`
	Installation *Installation        `json:"installation,omitempty"`
}

func (pr *PushEvent) GetInstallation() *Installation {
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
	Action       *string       `json:"action,omitempty"`
	Repositories []*Repository `json:"repositories,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
	// TODO key "requester" is not covered
}

func (ie *InstallationEvent) GetInstallation() *Installation {
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

type User struct {
	Login *string `json:"login,omitempty"`
	Email *string `json:"email,omitempty"`
}

func (u *User) GetLogin() string {
	if u == nil || u.Login == nil {
		return ""
	}
	return *u.Login
}

func (u *User) GetEmail() string {
	if u == nil || u.Email == nil {
		return ""
	}
	return *u.Email
}

type InstallationRepositoriesEvent struct {
	RepositoriesAdded   []*Repository `json:"repositories_added,omitempty"`
	RepositoriesRemoved []*Repository `json:"repositories_removed,omitempty"`
	RepositorySelection *string       `json:"repository_selection,omitempty"`
	Installation        *Installation `json:"installation,omitempty"`
}

func (ire *InstallationRepositoriesEvent) GetInstallation() *Installation {
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

type RepositoryCommit struct {
	SHA *string `json:"sha,omitempty"`
}

func (rc *RepositoryCommit) GetSHA() string {
	if rc == nil || rc.SHA == nil {
		return ""
	}
	return *rc.SHA
}

type StatusEvent struct {
	Context     *string `json:"context,omitempty"`
	State       *string `json:"state,omitempty"`
	Description *string `json:"description,omitempty"`
	TargetURL   *string `json:"target_url,omitempty"`

	Repo         *Repository       `json:"repository,omitempty"`
	Installation *Installation     `json:"installation,omitempty"`
	Commit       *RepositoryCommit `json:"commit,omitempty"`
	UpdatedAt    *Timestamp        `json:"updated_at,omitempty"`
	CreatedAt    *Timestamp        `json:"created_at,omitempty"`
}

func (se *StatusEvent) GetContext() string {
	if se == nil || se.Context == nil {
		return ""
	}
	return *se.Context
}

func (se *StatusEvent) GetCreatedAt() Timestamp {
	if se == nil || se.CreatedAt == nil {
		return Timestamp{}
	}
	return *se.CreatedAt
}

func (se *StatusEvent) GetUpdatedAt() Timestamp {
	if se == nil || se.UpdatedAt == nil {
		return Timestamp{}
	}
	return *se.UpdatedAt
}

func (se *StatusEvent) GetState() string {
	if se == nil || se.State == nil {
		return ""
	}
	return *se.State
}

func (se *StatusEvent) GetCommit() *RepositoryCommit {
	if se == nil {
		return nil
	}
	return se.Commit
}

func (se *StatusEvent) GetInstallation() *Installation {
	if se == nil {
		return nil
	}
	return se.Installation
}

func (se *StatusEvent) GetRepo() *Repository {
	if se == nil {
		return nil
	}
	return se.Repo
}

type WorkflowJobEvent struct {
	WorkflowJob  *WorkflowJob  `json:"workflow_job,omitempty"`
	Repo         *Repository   `json:"repository,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

func (wje *WorkflowJobEvent) GetInstallation() *Installation {
	if wje == nil {
		return nil
	}
	return wje.Installation
}

func (wje *WorkflowJobEvent) GetRepo() *Repository {
	if wje == nil {
		return nil
	}
	return wje.Repo
}

func (wje *WorkflowJobEvent) GetWorkflowJob() *WorkflowJob {
	if wje == nil {
		return nil
	}
	return wje.WorkflowJob
}

type WorkflowJob struct {
	HeadSHA     *string    `json:"head_sha,omitempty"`
	Status      *string    `json:"status,omitempty"`
	Conclusion  *string    `json:"conclusion,omitempty"`
	StartedAt   *Timestamp `json:"started_at,omitempty"`
	CompletedAt *Timestamp `json:"completed_at,omitempty"`
	Name        *string    `json:"name,omitempty"`
	HTMLURL     *string    `json:"html_url,omitempty"`
}

func (wj *WorkflowJob) GetStatus() string {
	if wj == nil || wj.Status == nil {
		return ""
	}
	return *wj.Status
}

func (wj *WorkflowJob) GetConclusion() string {
	if wj == nil || wj.Conclusion == nil {
		return ""
	}
	return *wj.Conclusion
}

func (wj *WorkflowJob) GetHeadSHA() string {
	if wj == nil || wj.HeadSHA == nil {
		return ""
	}
	return *wj.HeadSHA
}

func (wj *WorkflowJob) GetName() string {
	if wj == nil || wj.Name == nil {
		return ""
	}
	return *wj.Name
}

func (wj *WorkflowJob) GetCompletedAt() Timestamp {
	if wj == nil || wj.StartedAt == nil {
		return Timestamp{}
	}
	return *wj.CompletedAt
}

func (wj *WorkflowJob) GetStartedAt() Timestamp {
	if wj == nil || wj.StartedAt == nil {
		return Timestamp{}
	}
	return *wj.StartedAt
}

// Timestamp represents a time that can be unmarshalled from a JSON string
// formatted as either an RFC3339 or Unix timestamp. This is necessary for some
// fields since the GitHub API is inconsistent in how it represents times. All
// exported methods of time.Time can be called on Timestamp.
type Timestamp struct {
	time.Time
}

func (t Timestamp) String() string {
	return t.Time.String()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// Time is expected in RFC3339 or Unix format.
func (t *Timestamp) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	i, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		t.Time = time.Unix(i, 0)
		if t.Time.Year() > 3000 {
			t.Time = time.Unix(0, i*1e6)
		}
	} else {
		t.Time, err = time.Parse(`"`+time.RFC3339+`"`, str)
	}
	return
}

// Equal reports whether t and u are equal based on time.Equal
func (t Timestamp) Equal(u Timestamp) bool {
	return t.Time.Equal(u.Time)
}

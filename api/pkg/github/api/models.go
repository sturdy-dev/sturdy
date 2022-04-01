package api

import (
	"strconv"
	"time"

	"github.com/google/go-github/v39/github"
)

// this is copied from: https://github.com/google/go-github/blob/41cbf81925093b28eed72d7ef81577a7c6fdfd1b/github
// if we need more fields, or models, they should be copied manually to have control over the github's api response size

type PullRequest struct {
	ID             *int64     `json:"id,omitempty"`
	Title          *string    `json:"title,omitempty"`
	Number         *int       `json:"number,omitempty"`
	State          *string    `json:"state,omitempty"`
	Merged         *bool      `json:"merged,omitempty"`
	MergeCommitSHA *string    `json:"merge_commit_sha,omitempty"`
	HTMLURL        *string    `json:"html_url,omitempty"`
	ClosedAt       *time.Time `json:"closed_at,omitempty"`
	MergedAt       *time.Time `json:"merged_at,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	Body           *string    `json:"body,omitempty"`
	User           *User      `json:"user,omitempty"`

	Head *PullRequestBranch `json:"head,omitempty"`
	Base *PullRequestBranch `json:"base,omitempty"`
}

func ConvertPullRequest(pr *github.PullRequest) *PullRequest {
	return &PullRequest{
		ID:             pr.ID,
		Title:          pr.Title,
		Number:         pr.Number,
		State:          pr.State,
		Merged:         pr.Merged,
		MergeCommitSHA: pr.MergeCommitSHA,
		HTMLURL:        pr.HTMLURL,
		ClosedAt:       pr.ClosedAt,
		MergedAt:       pr.MergedAt,
		CreatedAt:      pr.CreatedAt,
		Body:           pr.Body,
		User:           ConvertUser(pr.User),
		Head:           ConvertPullReqestBranch(pr.Head),
		Base:           ConvertPullReqestBranch(pr.Base),
	}
}
func (pr *PullRequest) GetUser() *User {
	if pr == nil {
		return nil
	}
	return pr.User
}

func (pr *PullRequest) GetTitle() string {
	if pr.Title != nil {
		return *pr.Title
	}
	return ""
}

func (pr *PullRequest) GetNumber() int {
	if pr.Number != nil {
		return *pr.Number
	}
	return 0
}

func (pr *PullRequest) GetBody() string {
	if pr == nil || pr.Body == nil {
		return ""
	}
	return *pr.Body
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

func (pr *PullRequest) GetHead() *PullRequestBranch {
	if pr == nil {
		return nil
	}
	return pr.Head
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

func (pr *PullRequest) GetCreatedAt() time.Time {
	if pr == nil || pr.CreatedAt == nil {
		return time.Time{}
	}
	return *pr.CreatedAt
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

func ConvertPullReqestBranch(branch *github.PullRequestBranch) *PullRequestBranch {
	return &PullRequestBranch{
		Label: branch.Label,
		Ref:   branch.Ref,
		SHA:   branch.SHA,
		Repo:  ConvertRepository(branch.Repo),
		User:  ConvertUser(branch.User),
	}
}

func (prb *PullRequestBranch) GetUser() *User {
	if prb == nil {
		return nil
	}
	return prb.User
}

func (prb *PullRequestBranch) GetRef() string {
	if prb == nil || prb.Ref == nil {
		return ""
	}
	return *prb.Ref
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

func ConvertRepository(repo *github.Repository) *Repository {
	return &Repository{
		ID:       repo.ID,
		Name:     repo.Name,
		FullName: repo.FullName,
	}
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

type User struct {
	ID    *int64  `json:"id,omitempty"`
	Login *string `json:"login,omitempty"`
	Email *string `json:"email,omitempty"`
	Name  *string `json:"name,omitempty"`
}

func ConvertUser(user *github.User) *User {
	return &User{
		ID:    user.ID,
		Login: user.Login,
		Email: user.Email,
	}
}

func (u *User) GetID() int64 {
	if u == nil || u.ID == nil {
		return 0
	}
	return *u.ID
}

func (u *User) GetLogin() string {
	if u == nil || u.Login == nil {
		return ""
	}
	return *u.Login
}

func (u *User) GetName() string {
	if u == nil || u.Name == nil {
		return ""
	}
	return *u.Name
}

func (u *User) GetEmail() string {
	if u == nil || u.Email == nil {
		return ""
	}
	return *u.Email
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

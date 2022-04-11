package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"getsturdy.com/api/pkg/github"
)

func Test_primaryPullRequest(t *testing.T) {
	date := func(in string) time.Time {
		ts, err := time.Parse("2006-01-02", in)
		assert.NoError(t, err)
		return ts
	}

	oldNonFork := &github.PullRequest{
		ID:        "old-non-fork",
		CreatedAt: date("2020-05-05"),
		Fork:      false,
		State:     github.PullRequestStateOpen,
	}

	midFork := &github.PullRequest{
		ID:        "mid-fork",
		CreatedAt: date("2021-06-06"),
		Fork:      true,
		State:     github.PullRequestStateOpen,
	}

	newNonFork := &github.PullRequest{
		ID:        "new-non-fork",
		CreatedAt: date("2022-08-08"),
		Fork:      false,
		State:     github.PullRequestStateOpen,
	}

	openNonFork := &github.PullRequest{
		ID:        "open-non-fork",
		CreatedAt: date("2022-08-08"),
		Fork:      false,
		State:     github.PullRequestStateOpen,
	}

	closedNonFork := &github.PullRequest{
		ID:        "closed-non-fork",
		CreatedAt: date("2022-08-08"),
		Fork:      false,
		State:     github.PullRequestStateClosed,
	}

	tests := []struct {
		name    string
		arg     []*github.PullRequest
		want    *github.PullRequest
		wantErr error
	}{
		{name: "no-prs", arg: nil, want: nil, wantErr: ErrNotFound}, // xoxo
		{name: "newest-first", arg: []*github.PullRequest{oldNonFork, newNonFork}, want: newNonFork},
		{name: "newest-first-reverse", arg: []*github.PullRequest{newNonFork, oldNonFork}, want: newNonFork},
		{name: "prioritize-non-fork", arg: []*github.PullRequest{oldNonFork, midFork}, want: oldNonFork},
		{name: "prioritize-non-fork-all", arg: []*github.PullRequest{newNonFork, oldNonFork, midFork}, want: newNonFork},
		{name: "prioritize-open-over-closed", arg: []*github.PullRequest{openNonFork, closedNonFork}, want: openNonFork},
		{name: "prioritize-open-over-fork", arg: []*github.PullRequest{midFork, oldNonFork, openNonFork, closedNonFork}, want: openNonFork},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := primaryPullRequest(tt.arg)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equalf(t, tt.want, res, "primaryPullRequest(%v)", tt.arg)
		})
	}
}

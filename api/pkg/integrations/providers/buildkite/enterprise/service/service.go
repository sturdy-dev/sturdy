package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/integrations/providers"
	"getsturdy.com/api/pkg/integrations/providers/buildkite"
	db_buildkite "getsturdy.com/api/pkg/integrations/providers/buildkite/enterprise/db"
)

var _ providers.BuildProvider = &Service{}

type Service struct {
	configRepo db_buildkite.Repository
}

func New(configRepo db_buildkite.Repository) *Service {
	s := &Service{
		configRepo: configRepo,
	}
	providers.Register(s)
	return s
}

func (b *Service) ProviderType() providers.ProviderType {
	return providers.ProviderTypeBuild
}

func (b *Service) ProviderName() providers.ProviderName {
	return providers.ProviderNameBuildkite
}

func (b *Service) CreateIntegration(ctx context.Context, cfg *buildkite.Config) error {
	return b.configRepo.Create(ctx, cfg)
}

func (b *Service) UpdateIntegration(ctx context.Context, cfg *buildkite.Config) error {
	return b.configRepo.Update(ctx, cfg)
}

func (b *Service) GetConfigurationsByCodebaseID(ctx context.Context, codebaseID codebases.ID) ([]*buildkite.Config, error) {
	return b.configRepo.GetConfigsByCodebaseID(ctx, codebaseID)
}

func (b *Service) GetConfigurationByIntegrationID(ctx context.Context, integrationID string) (*buildkite.Config, error) {
	return b.configRepo.GetConfigByIntegrationID(ctx, integrationID)
}

func (b *Service) CreateBuild(ctx context.Context, integrationID, ciCommitId, title string) (*providers.Build, error) {
	cfg, err := b.configRepo.GetConfigByIntegrationID(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get config by codebase id: %w", err)
	}

	data := createBuild{
		Commit:  ciCommitId,
		Branch:  "main",
		Message: title,
		Author: createBuildAuthor{
			Name:  "Sturdy",
			Email: "no-reply@getsturdy.com",
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to build json: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url(cfg), bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create build request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+cfg.APIToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make create build request: %w", err)
	}
	defer resp.Body.Close()

	resContents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read contents: %w", err)
	}

	var parsedRes createBuildRes
	if err := json.Unmarshal(resContents, &parsedRes); err != nil {
		return nil, fmt.Errorf("failed to read response (%s): %w", string(resContents), err)
	}

	if parsedRes.ID == "" {
		return nil, fmt.Errorf("unexpected response, id not set")
	}

	return &providers.Build{
		Name: parsedRes.Pipeline.Name,
		URL:  parsedRes.WebURL,
	}, nil
}

func url(cfg *buildkite.Config) string {
	return fmt.Sprintf("https://api.buildkite.com/v2/organizations/%s/pipelines/%s/builds", slugify(cfg.OrganizationName), slugify(cfg.PipelineName))
}

// slugify converts a string to a kebeb-case slug
func slugify(str string) string {
	var buf bytes.Buffer
	for _, r := range str {
		if r == ' ' {
			buf.WriteRune('-')
		} else if r >= 'A' && r <= 'Z' {
			buf.WriteRune(r + 32)
		} else if r >= 'a' && r <= 'z' {
			buf.WriteRune(r)
		} else if r >= '0' && r <= '9' {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

type createBuild struct {
	Commit  string            `json:"commit"`
	Branch  string            `json:"branch"`
	Message string            `json:"message"`
	Author  createBuildAuthor `json:"author"`
}

type createBuildAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type createBuildRes struct {
	ID       string `json:"id"`
	WebURL   string `json:"web_url"`
	Number   int64  `json:"number"`
	Pipeline struct {
		Name string `json:"name"`
	} `json:"pipeline"`
}

/*
Example response from the API

{
  "id": "6a91d99b-9816-49dd-b804-17cfb6ac2ab5",
  "graphql_id": "QnVpbGQtLS02YTkxZDk5Yi05ODE2LTQ5ZGQtYjgwNC0xN2NmYjZhYzJhYjU=",
  "url": "https://api.buildkite.com/v2/organizations/sturdy/pipelines/sturdy/builds/2",
  "web_url": "https://buildkite.com/sturdy/sturdy/builds/2",
  "number": 2,
  "state": "scheduled",
  "blocked": false,
  "blocked_state": "",
  "message": "Sturdy",
  "commit": "192616070aa397c1400cfe909ee323c1520f06fa",
  "branch": "main",
  "tag": null,
  "env": {},
  "source": "api",
  "author": {
    "name": "Gustav Westling (fake data!)",
    "email": "gustav@getsturdy.com"
  },
  "creator": {
    "id": "6b8585f6-5004-4bbe-bf88-2ed46a876baf",
    "graphql_id": "VXNlci0tLTZiODU4NWY2LTUwMDQtNGJiZS1iZjg4LTJlZDQ2YTg3NmJhZg==",
    "name": "Gustav Westling",
    "email": "gustav@getsturdy.com",
    "avatar_url": "https://www.gravatar.com/avatar/b591c9056e1fe64f4d4409962113e0f6",
    "created_at": "2021-10-12T11:01:18.426Z"
  },
  "created_at": "2021-10-14T12:12:15.898Z",
  "scheduled_at": "2021-10-14T12:12:15.802Z",
  "started_at": null,
  "finished_at": null,
  "meta_data": {},
  "pull_request": null,
  "rebuilt_from": null,
  "pipeline": {
    "id": "dacdc97a-4abb-4c13-97df-6c4c00a2c469",
    "graphql_id": "UGlwZWxpbmUtLS1kYWNkYzk3YS00YWJiLTRjMTMtOTdkZi02YzRjMDBhMmM0Njk=",
    "url": "https://api.buildkite.com/v2/organizations/sturdy/pipelines/sturdy",
    "web_url": "https://buildkite.com/sturdy/sturdy",
    "name": "Sturdy",
    "description": null,
    "slug": "sturdy",
    "repository": "git@git.getsturdy.com:foo/bar.git",
    "cluster_id": null,
    "branch_configuration": null,
    "default_branch": "main",
    "skip_queued_branch_builds": false,
    "skip_queued_branch_builds_filter": null,
    "cancel_running_branch_builds": false,
    "cancel_running_branch_builds_filter": null,
    "allow_rebuilds": true,
    "provider": {
      "id": "private",
      "settings": {}
    },
    "builds_url": "https://api.buildkite.com/v2/organizations/sturdy/pipelines/sturdy/builds",
    "badge_url": "https://badge.buildkite.com/da7ae37ad772c11c523903876a4e3b0434d5059f27c6e5f90d.svg",
    "created_by": {
      "id": "6b8585f6-5004-4bbe-bf88-2ed46a876baf",
      "graphql_id": "VXNlci0tLTZiODU4NWY2LTUwMDQtNGJiZS1iZjg4LTJlZDQ2YTg3NmJhZg==",
      "name": "Gustav Westling",
      "email": "gustav@getsturdy.com",
      "avatar_url": "https://www.gravatar.com/avatar/b591c9056e1fe64f4d4409962113e0f6",
      "created_at": "2021-10-12T11:01:18.426Z"
    },
    "created_at": "2021-10-12T11:03:26.260Z",
    "archived_at": null,
    "env": {},
    "scheduled_builds_count": 0,
    "running_builds_count": 0,
    "scheduled_jobs_count": 0,
    "running_jobs_count": 1,
    "waiting_jobs_count": 0,
    "visibility": "private",
    "tags": null,
    "steps": [
      {
        "type": "script",
        "name": "",
        "command": "echo \"Hello\"",
        "artifact_paths": "",
        "branch_configuration": "",
        "env": {},
        "timeout_in_minutes": null,
        "agent_query_rules": [],
        "concurrency": null,
        "parallelism": null
      }
    ]
  },
  "jobs": [
    {
      "id": "1457935b-337f-440c-b241-06d763dc0bb0",
      "graphql_id": "Sm9iLS0tMTQ1NzkzNWItMzM3Zi00NDBjLWIyNDEtMDZkNzYzZGMwYmIw",
      "type": "script",
      "name": "",
      "step_key": null,
      "priority": {
        "number": 0
      },
      "agent_query_rules": [],
      "state": "scheduled",
      "build_url": "https://api.buildkite.com/v2/organizations/sturdy/pipelines/sturdy/builds/2",
      "web_url": "https://buildkite.com/sturdy/sturdy/builds/2#1457935b-337f-440c-b241-06d763dc0bb0",
      "log_url": "https://api.buildkite.com/v2/organizations/sturdy/pipelines/sturdy/builds/2/jobs/1457935b-337f-440c-b241-06d763dc0bb0/log",
      "raw_log_url": "https://api.buildkite.com/v2/organizations/sturdy/pipelines/sturdy/builds/2/jobs/1457935b-337f-440c-b241-06d763dc0bb0/log.txt",
      "artifacts_url": "https://api.buildkite.com/v2/organizations/sturdy/pipelines/sturdy/builds/2/jobs/1457935b-337f-440c-b241-06d763dc0bb0/artifacts",
      "command": "echo \"Hello\"",
      "soft_failed": false,
      "exit_status": null,
      "artifact_paths": "",
      "agent": null,
      "created_at": "2021-10-14T12:12:15.880Z",
      "scheduled_at": "2021-10-14T12:12:15.880Z",
      "runnable_at": "2021-10-14T12:12:15.969Z",
      "started_at": null,
      "finished_at": null,
      "retried": false,
      "retried_in_job_id": null,
      "retries_count": null,
      "parallel_group_index": null,
      "parallel_group_total": null
    }
  ]
}
*/

ALTER TABLE github_repositories
    DROP COLUMN deleted_at;

DROP INDEX github_repositories_installation_id_github_repository_id_idx;
CREATE UNIQUE INDEX ON github_repositories (installation_id, github_repository_id);

DROP INDEX github_prs_github_id_idx;
CREATE UNIQUE INDEX ON github_pull_requests (github_id);
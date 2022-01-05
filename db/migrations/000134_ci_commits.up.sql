CREATE TABLE ci_commits
(
    id                TEXT PRIMARY KEY,
    codebase_id       TEXT                     NOT NULL,
    ci_repo_commit_id TEXT                     NOT NULL,
    trunk_commit_id   TEXT                     NOT NULL,
    created_at        TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX ci_commits_codebase_id_commit_id_idx
    ON ci_commits (codebase_id, ci_repo_commit_id);
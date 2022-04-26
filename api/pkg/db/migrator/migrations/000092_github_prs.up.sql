CREATE TABLE IF NOT EXISTS github_prs
(
    id                       text                     NOT NULL,
    workspace_id             text                     NOT NULL,
    github_id                int                      NOT NULL,
    github_repository_id     int                      NOT NULL,
    github_author_user_login text                     NOT NULL,
    github_pr_number         int                      NOT NULL,
    head                     text                     NOT NULL,
    base                     text                     NOT NULL,
    opened                   boolean                  NOT NULL,
    merged                   boolean                  NOT NULL,
    created_at               timestamp WITH TIME ZONE NOT NULL,
    updated_at               timestamp WITH TIME ZONE, -- Nullable
    closed_at                timestamp WITH TIME ZONE, -- Nullable
    merged_at                timestamp WITH TIME ZONE, -- Nullable
    PRIMARY KEY (id)
);

CREATE INDEX github_prs_workspace_id_idx
    ON github_prs (workspace_id);
CREATE UNIQUE INDEX ON github_prs (github_id);

CREATE TABLE IF NOT EXISTS github_installations
(
    id              text      NOT NULL,
    installation_id int       NOT NULL,
    owner           text      NOT NULL,
    created_at      timestamp NOT NULL,
    uninstalled_at  timestamp,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX ON github_installations (installation_id);

CREATE TABLE IF NOT EXISTS github_repositories
(
    id                                   text      NOT NULL,
    installation_id                      int       NOT NULL,
    name                                 text      NOT NULL,
    github_repository_id                 int       NOT NULL,
    created_at                           timestamp NOT NULL,
    uninstalled_at                       timestamp,
    synced_at                            timestamp,
    installation_access_token            text,
    installation_access_token_expires_at timestamp,
    codebase_id                          text      NOT NULL,
    tracked_branch                       text,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX ON github_repositories (codebase_id);
CREATE UNIQUE INDEX ON github_repositories (installation_id, github_repository_id);

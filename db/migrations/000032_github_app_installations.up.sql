CREATE TABLE IF NOT EXISTS github_app_installations
(
    installation_id int       NOT NULL,
    owner           text      NOT NULL,
    created_at      timestamp NOT NULL,
    PRIMARY KEY (installation_id)
);

CREATE TABLE IF NOT EXISTS github_app_installed_repositories
(
    installation_id                      int       NOT NULL,
    name                                 text      NOT NULL,
    github_repository_id                 int       NOT NULL,
    created_at                           timestamp NOT NULL,
    uninstalled_at                       timestamp,
    installation_access_token            text,
    installation_access_token_expires_at timestamp,
    repository_id                        text      NOT NULL,
    PRIMARY KEY (installation_id, github_repository_id)
);

CREATE UNIQUE INDEX ON github_app_installed_repositories (repository_id);
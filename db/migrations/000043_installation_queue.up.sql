ALTER TABLE github_app_installed_repositories ADD COLUMN synced_at TIMESTAMP;

CREATE INDEX github_app_installed_repositories_synced_at_uninstalled_atidx
    ON github_app_installed_repositories (synced_at, uninstalled_at);

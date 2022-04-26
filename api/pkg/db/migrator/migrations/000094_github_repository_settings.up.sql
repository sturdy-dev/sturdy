ALTER TABLE github_repositories
    ADD COLUMN integration_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN github_source_of_truth BOOLEAN NOT NULL DEFAULT TRUE;
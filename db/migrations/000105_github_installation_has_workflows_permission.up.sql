ALTER TABLE github_installations
    ADD COLUMN has_workflows_permission BOOLEAN NOT NULL DEFAULT FALSE;
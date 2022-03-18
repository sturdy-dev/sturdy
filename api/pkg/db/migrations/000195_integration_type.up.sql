ALTER TABLE ci_configurations
    ADD COLUMN provider_type TEXT NOT NULL DEFAULT 'build';

ALTER TABLE ci_configurations RENAME TO integrations;
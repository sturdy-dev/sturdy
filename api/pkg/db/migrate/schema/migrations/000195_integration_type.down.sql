ALTER TABLE integrations RENAME TO ci_configurations;

ALTER TABLE ci_configurations
    DROP COLUMN provider_type;
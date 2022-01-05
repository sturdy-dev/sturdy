ALTER TABLE ci_configurations_buildkite
    ADD  COLUMN id                TEXT                     NOT NULL PRIMARY KEY,
    ADD  COLUMN organization_name TEXT                     NOT NULL,
    ADD  COLUMN pipeline_name     TEXT                     NOT NULL,
    ADD  COLUMN created_at        TIMESTAMP WITH TIME ZONE NOT NULL,
    ADD  COLUMN updated_at        TIMESTAMP WITH TIME ZONE,
    DROP COLUMN organization_slug,
    DROP COLUMN pipeline_slug;

ALTER TABLE ci_configurations
    ALTER COLUMN seed_files DROP NOT NULL;

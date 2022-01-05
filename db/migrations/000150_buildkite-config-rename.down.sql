ALTER TABLE ci_configurations_buildkite
    DROP COLUMN id,
    DROP COLUMN organization_name,
    DROP COLUMN pipeline_name,
    DROP COLUMN created_at,
    DROP COLUMN updated_at,
    ADD  COLUMN organization_slug TEXT NOT NULL,
    ADD  COLUMN pipeline_slug     TEXT NOT NULL;

ALTER TABLE ci_configurations
    ALTER COLUMN seed_files SET NOT NULL;

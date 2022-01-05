TRUNCATE ci_configurations_buildkite;

ALTER TABLE ci_configurations_buildkite
    ADD COLUMN integration_id TEXT NOT NULL UNIQUE;

ALTER TABLE ci_configurations DROP CONSTRAINT ci_configurations_codebase_id_provider_key;

ALTER TABLE ci_configurations_buildkite DROP CONSTRAINT ci_configurations_buildkite_codebase_id_key;
CREATE TABLE ci_configurations_buildkite (
    codebase_id   TEXT NOT NULL,
    org_slug      TEXT NOT NULL,
    pipeline_slug TEXT NOT NULL,
    token         TEXT NOT NULL,

    UNIQUE(codebase_id)
);

CREATE TABLE ci_configurations (
    codebase_id TEXT      NOT NULL,
    provider    TEXT      NOT NULL,
    seed_files  TEXT[]    NOT NULL,
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL,

    UNIQUE(codebase_id, provider)
);

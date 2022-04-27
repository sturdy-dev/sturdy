CREATE TABLE statuses (
    id          TEXT      NOT NULL PRIMARY KEY,
    codebase_id TEXT      NOT NULL,
    commit_id   TEXT      NOT NULL,
    title       TEXT      NOT NULL,
    description TEXT,
    type        TEXT      NOT NULL,
    timestamp   TIMESTAMP NOT NULL
);

CREATE INDEX statuses_commit_id_codebase_id_ix
    ON statuses (commit_id, codebase_id);

CREATE INDEX statuses_commit_id_codebase_id_title_timestamp_ix ON
    statuses (commit_id, codebase_id, title, timestamp DESC);

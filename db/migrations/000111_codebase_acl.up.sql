CREATE TABLE acls
(
    id          TEXT                     PRIMARY KEY,
    codebase_id TEXT                     NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    policy      TEXT                     NOT NULL
);

CREATE INDEX acls_codebase_id_idx ON
    acls (codebase_id);

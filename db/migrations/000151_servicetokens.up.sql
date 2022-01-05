CREATE TABLE servicetokens (
    id            TEXT                     NOT NULL PRIMARY KEY,
    codebase_id   TEXT                     NOT NULL,
    hash          BYTEA                    NOT NULL,
    name          TEXT                     NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL,
    last_used_at  TIMESTAMP WITH TIME ZONE
)

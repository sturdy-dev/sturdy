START TRANSACTION;

CREATE TABLE change_commits
(
    change_id   TEXT NOT NULL,
    commit_id   TEXT NOT NULL,
    codebase_id TEXT NOT NULL
);

-- NON UNIQUE INDEX for change_id
CREATE INDEX change_commits_change_id_idx
    ON change_commits (change_id);

-- UNIQUE INDEX for commit_id and codebase_id
CREATE UNIQUE INDEX change_commits_commit_id_codebase_id_idx
    ON change_commits (commit_id, codebase_id);

-- Migrate data to new table
INSERT INTO change_commits (change_id, codebase_id, commit_id) (SELECT id AS change_id, codebase_id, commit_id FROM changes);

-- Drop column from previous table
ALTER TABLE changes DROP COLUMN commit_id;

COMMIT
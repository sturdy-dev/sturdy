DROP INDEX changes_codebase_id_commit_id_idx;

CREATE INDEX changes_codebase_id_commit_id_idx
    ON changes (codebase_id, commit_id);


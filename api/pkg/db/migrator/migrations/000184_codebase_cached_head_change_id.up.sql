ALTER TABLE codebases
    ADD COLUMN calculated_head_change_id BOOL NOT NULL DEFAULT FALSE,
    ADD COLUMN cached_head_change_id TEXT;

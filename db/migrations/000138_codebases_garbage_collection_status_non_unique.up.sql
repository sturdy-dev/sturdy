ALTER TABLE codebases_garbage_collection_status
    DROP CONSTRAINT codebases_garbage_collection_status_pkey;

ALTER TABLE codebases_garbage_collection_status
    RENAME COLUMN last_collected_completed_at TO completed_at;

ALTER TABLE codebases_garbage_collection_status
    RENAME COLUMN last_collection_duration_millis TO duration_millis;

ALTER TABLE codebases_garbage_collection_status
    SET codebase_id PRIMARY KEY;


ALTER TABLE codebases_garbage_collection_status
    RENAME COLUMN completed_at TO last_collected_completed_at;

ALTER TABLE codebases_garbage_collection_status
    RENAME COLUMN duration_millis TO last_collection_duration_millis;

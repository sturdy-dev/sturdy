ALTER TABLE conflicting_check_log ADD COLUMN working_directory_hash TEXT;

CREATE INDEX conflicting_check_log_repo_base_onto_working_idx
    ON conflicting_check_log (repository_id, base, onto, working_directory_hash);

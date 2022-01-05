ALTER TABLE conflicting_check_log ADD COLUMN is_conflict_in_working_directory BOOLEAN;
UPDATE conflicting_check_log SET is_conflict_in_working_directory = false;
ALTER TABLE conflicting_check_log ALTER COLUMN is_conflict_in_working_directory SET NOT NULL;

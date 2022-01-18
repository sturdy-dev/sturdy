ALTER TABLE conflict_repo ADD COLUMN fetched_at timestamp ;
UPDATE conflict_repo SET fetched_at = created_at;
ALTER TABLE conflict_repo ALTER COLUMN fetched_at SET NOT NULL;

-- backfill short_id for the codebases that haven't been automatically migrated yet
UPDATE codebases SET short_id = substring(id, 0, 8)  WHERE short_id IS NULL;

ALTER TABLE codebases
    ALTER COLUMN short_id SET NOT NULL;
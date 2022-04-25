DELETE FROM snapshots
WHERE workspace_id IS NULL;

ALTER TABLE snapshots
    ALTER COLUMN workspace_id SET NOT NULL;

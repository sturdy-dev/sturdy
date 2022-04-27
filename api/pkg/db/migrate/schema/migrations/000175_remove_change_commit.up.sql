DELETE FROM changes WHERE commit_id IS NULL;

DROP TABLE change_commits;

ALTER TABLE changes
    ALTER COLUMN commit_id SET NOT NULL;

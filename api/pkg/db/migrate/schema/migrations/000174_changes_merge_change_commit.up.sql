BEGIN;

ALTER TABLE changes
    ADD COLUMN commit_id TEXT;

WITH subquery AS (
    SELECT change_id, commit_id
    FROM change_commits
    WHERE trunk IS TRUE
)
UPDATE changes
SET commit_id = subquery.commit_id FROM subquery
WHERE changes.id = subquery.change_id;

-- VERIFY
-- SELECT count(*) FROM changes where commit_id is null;

-- TODO(gustav): in a future migration, complete the cleanup

-- DELETE FROM changes WHERE commit_id IS NULL;
-- DROP TABLE change_commits;
--
-- ALTER TABLE changes
--     ALTER COLUMN commit_id ADD NOT NULL;


COMMIT;

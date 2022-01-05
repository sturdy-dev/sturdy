START TRANSACTION;

DELETE FROM comments WHERE change_id IS NULL;

ALTER TABLE comments
    ALTER COLUMN change_id SET NOT NULL,
    DROP COLUMN workspace_id,
    DROP COLUMN context,
    DROP COLUMN context_starts_at_line;

COMMIT;
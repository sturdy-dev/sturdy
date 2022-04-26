ALTER TABLE changes
    ADD COLUMN num_comments INTEGER NOT NULL DEFAULT 0;

UPDATE changes SET num_comments = (
    SELECT COUNT(1) FROM comments WHERE comments.change_id = changes.id
);

ALTER TABLE comments
    ADD COLUMN parent_comment_id TEXT;

CREATE INDEX comments_parent_comment_id_idx
    ON comments(parent_comment_id);
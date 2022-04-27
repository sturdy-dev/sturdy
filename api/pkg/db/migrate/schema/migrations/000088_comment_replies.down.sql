ALTER TABLE comments
    DROP COLUMN parent_comment_id;

DELETE INDEX comments_parent_comment_id_idx;
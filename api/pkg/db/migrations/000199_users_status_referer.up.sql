ALTER TABLE users
    ADD COLUMN referer TEXT,
    ADD COLUMN status TEXT;

UPDATE users
    SET status = 'active';

ALTER TABLE users
    ALTER COLUMN status SET NOT NULL;

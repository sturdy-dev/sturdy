ALTER TABLE jwt_keys
    ADD COLUMN created_at TIMESTAMP;

UPDATE jwt_keys SET created_at = NOW();

ALTER TABLE jwt_keys
    ALTER COLUMN created_at SET DEFAULT NOW();

ALTER TABLE jwt_keys
    ALTER COLUMN created_at SET NOT NULL;

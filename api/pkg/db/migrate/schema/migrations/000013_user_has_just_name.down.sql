ALTER TABLE users
    RENAME COLUMN name TO first_name;
ALTER TABLE users
    ADD COLUMN last_name text;

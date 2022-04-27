ALTER TABLE users
    ADD COLUMN default_view text;
ALTER TABLE users
    ADD CONSTRAINT fk_view_id
        FOREIGN KEY (default_view)
            REFERENCES views (id);

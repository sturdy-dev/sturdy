ALTER TABLE codebases
    ADD COLUMN short_id TEXT;

CREATE UNIQUE INDEX codebases_short_id_idx
    ON codebases(short_id);
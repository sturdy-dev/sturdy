ALTER TABLE codebases
    ADD COLUMN invite_code TEXT;

CREATE UNIQUE INDEX idx_codebases_invite_code
    ON codebases (invite_code);
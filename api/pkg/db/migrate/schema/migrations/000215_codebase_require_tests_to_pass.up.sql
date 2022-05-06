ALTER TABLE codebases
    ADD COLUMN require_healthy_status BOOLEAN NOT NULL DEFAULT false;
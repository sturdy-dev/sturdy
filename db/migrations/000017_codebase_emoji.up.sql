ALTER TABLE codebases
    ADD COLUMN emoji text;

UPDATE codebases SET emoji = '🌟' WHERE emoji IS NULL;
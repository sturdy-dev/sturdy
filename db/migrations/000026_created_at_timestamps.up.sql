ALTER TABLE users ADD COLUMN created_at timestamp;
ALTER TABLE views ADD COLUMN created_at timestamp;
ALTER TABLE workspaces ADD COLUMN created_at timestamp;
ALTER TABLE codebases ADD COLUMN created_at timestamp;
ALTER TABLE codebase_users ADD COLUMN created_at timestamp;

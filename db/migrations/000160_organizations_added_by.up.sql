ALTER TABLE organizations
    ADD COLUMN created_by TEXT NOT NULL,
    ADD COLUMN deleted_by TEXT;

ALTER TABLE organization_members
    ADD COLUMN created_by TEXT NOT NULL,
    ADD COLUMN deleted_by TEXT;
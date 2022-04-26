ALTER TABLE organizations
    DROP COLUMN created_by,
    DROP COLUMN deleted_by;

ALTER TABLE organization_members
    DROP COLUMN created_by,
    DROP COLUMN deleted_by;
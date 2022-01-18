CREATE TABLE organizations
(
    id         TEXT                     NOT NULL PRIMARY KEY,
    name       TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE organization_members
(
    id              TEXT                     NOT NULL PRIMARY KEY,
    user_id         TEXT                     NOT NULL,
    organization_id TEXT                     NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at      TIMESTAMP WITH TIME ZONE
);

CREATE INDEX organization_members_user_id_idx ON organization_members (user_id);
CREATE INDEX organization_members_organization_id_idx ON organization_members (organization_id);
CREATE INDEX organization_members_user_id_organization_id_idx ON organization_members (user_id, organization_id);

ALTER TABLE codebases ADD COLUMN organization_id TEXT;
CREATE INDEX codebases_organization_id_idx ON codebases(organization_id);
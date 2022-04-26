DROP TABLE self_hosted_licenses;
DROP TABLE self_hosted_license_validations;

CREATE TABLE licenses (
    id              TEXT                     NOT NULL PRIMARY KEY,
    organization_id TEXT                     NOT NULL,
    key             TEXT                     NOT NULL UNIQUE,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL,
    expires_at      TIMESTAMP WITH TIME ZONE NOT NULL
);

ALTER TABLE installations
    ADD COLUMN license_key TEXT;

CREATE TABLE license_validations (
    id         TEXT                     NOT NULL PRIMARY KEY,
    license_id TEXT                     NOT NULL,
    timestamp  TIMESTAMP WITH TIME ZONE NOT NULL,
    status     TEXT                     NOT NULL
);

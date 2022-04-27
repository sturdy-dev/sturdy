CREATE TABLE self_hosted_licenses (
    id TEXT PRIMARY KEY,
    cloud_organization_id TEXT NOT NULL,
    seats INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    active BOOLEAN NOT NULL
);

CREATE TABLE self_hosted_license_validations (
    id TEXT PRIMARY KEY,
    self_hosted_license_id TEXT NOT NULL,
    validated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    status BOOLEAN NOT NULL,
    reported_version  TEXT NOT NULL,
    reported_booted_at TIMESTAMP WITH TIME ZONE NOT NULL,
    reported_user_count INT NOT NULL,
    reported_codebase_count INT NOT NULL,
    from_ip_addr TEXT NOT NULL
);

CREATE INDEX self_hosted_licenses_cloud_organization_id_idx ON self_hosted_licenses(cloud_organization_id);
CREATE INDEX self_hosted_license_validations_self_hosted_license_id_idx ON self_hosted_license_validations(self_hosted_license_id);

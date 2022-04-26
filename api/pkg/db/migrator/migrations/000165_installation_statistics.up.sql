CREATE TABLE installation_statistics (
    installation_id TEXT                     NOT NULL,
    license_key     TEXT                     NOT NULL,
    version         TEXT                     NOT NULL,
    ip              TEXT                     NOT NULL,
    recorded_at     TIMESTAMP WITH TIME ZONE NOT NULL,
    received_at     TIMESTAMP WITH TIME ZONE NOT NULL
);
CREATE INDEX installation_statistics_license_key_idx ON installation_statistics (installation_id);

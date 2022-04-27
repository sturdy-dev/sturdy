CREATE TABLE keypairs
(
    id           TEXT PRIMARY KEY,
    public_key   TEXT                     NOT NULL,
    private_key  TEXT                     NOT NULL,
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL,
    created_by   TEXT                     NOT NULL,
    last_used_at TIMESTAMP WITH TIME ZONE NOT NULL
);
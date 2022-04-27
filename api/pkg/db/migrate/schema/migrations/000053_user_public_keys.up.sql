CREATE TABLE user_public_keys (
    public_key TEXT NOT NULL,
    user_id TEXT NOT NULL,
    added_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP,
    PRIMARY KEY(public_key)
)
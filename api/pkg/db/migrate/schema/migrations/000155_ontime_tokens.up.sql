CREATE TABLE onetime_tokens (
    key        TEXT                     NOT NULL,
    user_id    TEXT                     NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    clicks     INTEGER                  NOT NULL
);

CREATE UNIQUE INDEX onetime_tokens_key_user_id_idx ON onetime_tokens (key, user_id) WHERE clicks = 0;

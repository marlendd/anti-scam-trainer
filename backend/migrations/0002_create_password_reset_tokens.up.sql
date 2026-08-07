CREATE TABLE password_reset_tokens (
                                       id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                       user_id    UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
                                       token_hash TEXT NOT NULL,
                                       expires_at TIMESTAMPTZ NOT NULL,
                                       used_at    TIMESTAMPTZ,
                                       created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX password_reset_tokens_token_hash_idx ON password_reset_tokens (token_hash);
CREATE INDEX password_reset_tokens_user_id_idx ON password_reset_tokens (user_id);
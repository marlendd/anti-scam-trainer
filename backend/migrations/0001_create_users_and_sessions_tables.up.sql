CREATE TABLE users (
                       id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       email         TEXT NOT NULL,
                       password_hash TEXT NOT NULL,
                       created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
                       updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX users_email_idx ON users (lower(email));

CREATE TABLE sessions (
                          id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                          user_id    UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
                          user_agent TEXT,
                          ip         TEXT,
                          expires_at TIMESTAMPTZ NOT NULL,
                          created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX sessions_user_id_idx ON sessions (user_id);
CREATE INDEX sessions_expires_at_idx ON sessions (expires_at);
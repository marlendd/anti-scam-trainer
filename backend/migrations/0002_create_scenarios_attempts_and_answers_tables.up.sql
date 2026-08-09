CREATE TABLE scenario_versions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    logical_id  UUID NOT NULL,
    version     INTEGER NOT NULL CHECK (version > 0),
    role        TEXT NOT NULL CHECK (role IN ('buyer', 'seller')),
    title       TEXT NOT NULL CHECK (btrim(title) <> ''),
    description TEXT NOT NULL,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    content     JSONB NOT NULL CHECK (jsonb_typeof(content) = 'object'),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (logical_id, version)
);

CREATE INDEX scenario_versions_role_active_idx
    ON scenario_versions (role, is_active, created_at DESC);

CREATE UNIQUE INDEX scenario_versions_one_active_idx
    ON scenario_versions (logical_id)
    WHERE is_active;

CREATE TABLE attempts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    scenario_id     UUID NOT NULL REFERENCES scenario_versions (id) ON DELETE RESTRICT,
    status          TEXT NOT NULL CHECK (status IN ('in_progress', 'completed', 'aborted')),
    current_node_id TEXT,
    ending_id       TEXT,
    score           INTEGER CHECK (score BETWEEN 0 AND 100),
    started_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at    TIMESTAMPTZ,
    CONSTRAINT attempts_state_check CHECK (
        (
            status = 'in_progress'
            AND current_node_id IS NOT NULL
            AND ending_id IS NULL
            AND score IS NULL
            AND completed_at IS NULL
        )
        OR
        (
            status = 'completed'
            AND current_node_id IS NULL
            AND ending_id IS NOT NULL
            AND score IS NOT NULL
            AND completed_at IS NOT NULL
        )
        OR
        (
            status = 'aborted'
            AND ending_id IS NULL
            AND score IS NULL
            AND completed_at IS NULL
        )
    )
);

CREATE INDEX attempts_user_started_idx
    ON attempts (user_id, started_at DESC);

CREATE INDEX attempts_user_status_idx
    ON attempts (user_id, status);

CREATE INDEX attempts_scenario_score_idx
    ON attempts (scenario_id, score DESC)
    WHERE status = 'completed';

CREATE TABLE answers (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attempt_id        UUID NOT NULL REFERENCES attempts (id) ON DELETE CASCADE,
    node_id           TEXT NOT NULL CHECK (btrim(node_id) <> ''),
    choice_id         TEXT NOT NULL CHECK (btrim(choice_id) <> ''),
    idempotency_key   UUID NOT NULL,
    weight            SMALLINT NOT NULL CHECK (weight BETWEEN 1 AND 3),
    choice_score      SMALLINT NOT NULL CHECK (choice_score IN (0, 50, 100)),
    risk_categories   JSONB NOT NULL DEFAULT '[]'::JSONB
        CHECK (jsonb_typeof(risk_categories) = 'array'),
    consequence       TEXT NOT NULL CHECK (btrim(consequence) <> ''),
    explanation       TEXT NOT NULL CHECK (btrim(explanation) <> ''),
    response          JSONB NOT NULL CHECK (jsonb_typeof(response) = 'object'),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (attempt_id, node_id),
    UNIQUE (idempotency_key),
    CONSTRAINT answers_risk_required_check CHECK (
        choice_score = 100 OR jsonb_array_length(risk_categories) > 0
    )
);

CREATE INDEX answers_attempt_created_idx
    ON answers (attempt_id, created_at);

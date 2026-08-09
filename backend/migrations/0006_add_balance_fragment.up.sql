ALTER TABLE attempts 
ADD COLUMN current_balance INTEGER NOT NULL DEFAULT 0;

ALTER TABLE answers 
ADD COLUMN balance_delta INTEGER NOT NULL DEFAULT 0;

ALTER TABLE scenario_versions 
ADD COLUMN reward_fragment_id TEXT
    CHECK (reward_fragment_id IS NULL OR btrim(reward_fragment_id) <> '');

CREATE TABLE user_inventory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    scenario_id UUID NOT NULL REFERENCES scenario_versions(id) ON DELETE RESTRICT,
    fragment_id TEXT NOT NULL CHECK (btrim(fragment_id) <> ''),
    earned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, fragment_id),
    UNIQUE(user_id, scenario_id)
);

CREATE INDEX user_inventory_user_earned_idx
    ON user_inventory (user_id, earned_at);

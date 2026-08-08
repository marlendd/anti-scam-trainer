ALTER TABLE attempts 
ADD COLUMN current_balance INTEGER NOT NULL DEFAULT 0;

ALTER TABLE answers 
ADD COLUMN balance_delta INTEGER NOT NULL DEFAULT 0;

ALTER TABLE scenario_versions 
ADD COLUMN reward_fragment_id TEXT; 

CREATE TABLE user_inventory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    scenario_id UUID NOT NULL REFERENCES scenario_versions(id), 
    fragment_id TEXT NOT NULL,                                  
    earned_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE(user_id, fragment_id) -- один фрагмент нельзя получить дважды
);
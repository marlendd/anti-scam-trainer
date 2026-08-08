DROP TABLE IF EXISTS user_inventory;

ALTER TABLE scenario_versions DROP COLUMN IF EXISTS reward_fragment_id;
ALTER TABLE answers DROP COLUMN IF EXISTS balance_delta;
ALTER TABLE attempts DROP COLUMN IF EXISTS current_balance;
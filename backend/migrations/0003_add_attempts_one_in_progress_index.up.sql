CREATE UNIQUE INDEX attempts_one_in_progress_idx
    ON attempts (user_id, scenario_id)
    WHERE status = 'in_progress';

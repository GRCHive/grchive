ALTER TABLE script_runs
ADD COLUMN build_success BOOLEAN DEFAULT false;

ALTER TABLE script_runs
ADD COLUMN run_success BOOLEAN DEFAULT false;

ALTER TABLE script_runs
DROP COLUMN end_time;

ALTER TABLE script_runs
ADD COLUMN build_finish_time TIMESTAMPTZ;

ALTER TABLE script_runs
ADD COLUMN run_finish_time TIMESTAMPTZ;

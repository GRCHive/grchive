ALTER TABLE script_runs
ADD COLUMN build_start_time TIMESTAMPTZ;

ALTER TABLE script_runs
ADD COLUMN run_start_time TIMESTAMPTZ;

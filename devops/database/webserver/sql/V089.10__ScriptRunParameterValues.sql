ALTER TABLE client_script_code_parameters
ADD COLUMN id BIGSERIAL PRIMARY KEY;

CREATE TABLE script_run_parameters (
    run_id BIGINT NOT NULL REFERENCES script_runs(id) ON DELETE CASCADE,
    param_id BIGINT NOT NULL REFERENCES client_script_code_parameters(id) ON DELETE CASCADE,
    vals JSONB NOT NULL,
    UNIQUE(run_id, param_id)
);

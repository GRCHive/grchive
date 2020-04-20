ALTER TABLE script_run_parameters
ADD COLUMN param_name VARCHAR(64);

UPDATE script_run_parameters
SET param_name = sub.name
FROM (
    SELECT *
    FROM client_script_code_parameters
) AS sub
WHERE sub.id = param_id;

ALTER TABLE script_run_parameters
DROP COLUMN param_id CASCADE;

CREATE UNIQUE INDEX script_run_parameters_run_id_param_name_idx ON script_run_parameters(run_id, param_name);

ALTER TABLE script_run_parameters
ADD CONSTRAINT run_id_param_name_unique
UNIQUE USING INDEX script_run_parameters_run_id_param_name_idx;

ALTER TABLE client_script_code_parameters
DROP COLUMN id;

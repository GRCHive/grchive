ALTER TABLE code_to_client_data_link 
DROP CONSTRAINT code_to_client_data_link_code_id_org_id_fkey,
ADD CONSTRAINT code_to_client_data_link_code_id_org_id_fkey
    FOREIGN KEY(code_id, org_id)
    REFERENCES managed_code(id, org_id)
    ON DELETE CASCADE;

ALTER TABLE code_to_client_data_link 
DROP CONSTRAINT code_to_client_data_link_data_id_org_id_fkey,
ADD CONSTRAINT code_to_client_data_link_data_id_org_id_fkey
    FOREIGN KEY(data_id, org_id)
    REFERENCES client_data(id, org_id)
    ON DELETE CASCADE;

DROP INDEX file_folder_link_org_id_folder_id_file_id_idx;
CREATE UNIQUE INDEX ON file_folder_link(org_id, folder_id, file_id);

DROP INDEX control_folder_link_org_id_folder_id_control_id_idx;
CREATE UNIQUE INDEX ON control_folder_link(org_id, folder_id, control_id);

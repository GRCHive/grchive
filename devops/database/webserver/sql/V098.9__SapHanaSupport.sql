INSERT INTO supported_databases (id, name, has_sql_support)
VALUES (3, 'SAP HANA', true);

INSERT INTO data_source_options (name, kotlin_class)
VALUES ('Root.Database.SAP HANA', 'grchive.core.data.sources.databases.SAPHanaDataSource');

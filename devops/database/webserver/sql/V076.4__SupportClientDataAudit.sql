CREATE OR REPLACE FUNCTION audit_resource_type_to_human_name(r VARCHAR(64))
    RETURNS TEXT AS
$$
    DECLARE
        nm TEXT;
    BEGIN
        SELECT CASE
            WHEN r='database_connection_info'
                THEN 'Database Connection Info'
            WHEN r='database_resources'
                THEN 'Database'
            WHEN r='database_sql_query_requests'
                THEN 'SQL Query Request'
	        WHEN r='database_sql_metadata'
                THEN 'SQL Query Metadata'
	        WHEN r='database_sql_queries'
                THEN 'SQL Query Version'
            WHEN r='document_requests'
                THEN 'Document Request'
            WHEN r='file_metadata'
                THEN 'Documentation Metadata'
            WHEN r='file_storage'
                THEN 'Documentation'
            WHEN r='general_ledger_accounts'
                THEN 'GL Account'
            WHEN r='general_ledger_categories'
                THEN 'GL Category'
            WHEN r='infrastructure_servers'
                THEN 'GL Servers'
            WHEN r='process_flows'
                THEN 'Process Flow'
            WHEN r='process_flow_controls'
                THEN 'Control'
            WHEN r='process_flow_control_documentation_categories'
                THEN 'Documentation Category'
            WHEN r='process_flow_nodes'
                THEN 'Process Flow Node'
            WHEN r='process_flow_node_inputs'
                THEN 'Process Flow Node (Input)'
            WHEN r='process_flow_node_outputs'
                THEN 'Process Flow Node (Output)'
            WHEN r='process_flow_risks'
                THEN 'Risk'
            WHEN r='systems'
                THEN 'System'
            WHEN r='vendors'
                THEN 'Vendor'
            WHEN r='vendor_products'
                THEN 'Vendor Product'
            WHEN r='client_data'
                THEN 'Data Object'
        END INTO nm;
        RETURN nm;
    END;
$$ LANGUAGE plpgsql;

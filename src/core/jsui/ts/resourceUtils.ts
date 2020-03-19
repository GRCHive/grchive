export interface ResourceHandle {
    displayText: string
    // resourceUri will be null only if the resource has been deleted.
    resourceUri: string | null
}

export function resourceTypeToIcon(typ : string) : string {
    switch (typ) {
    case 'database_sql_query_requests':
        return 'mdi-database-search'
    case 'document_requests':
        return 'mdi-file-search'
    case 'process_flow_controls':
        return 'mdi-shield-lock-outline'
    }
    return 'mdi-alert-circle-outline'
}

export function standardizeResourceType(typ : string) : string {
    switch (typ) {
        case 'database_connection_info':
            return 'Database Connection Info'
        case 'database_resources':
            return 'Database'
        case 'database_sql_query_requests':
            return 'SQL Query Request'
	    case 'database_sql_metadata':
            return 'SQL Query Metadata'
	    case 'database_sql_queries':
            return 'SQL Query Version'
        case 'document_requests':
            return 'Document Request'
        case 'file_metadata':
            return 'Documentation Metadata'
        case 'file_storage':
            return 'Documentation'
        case 'general_ledger_accounts':
            return 'GL Account'
        case 'general_ledger_categories':
            return 'GL Category'
        case 'infrastructure_servers':
            return 'Servers'
        case 'process_flows':
            return 'Process Flow'
        case 'process_flow_controls':
            return 'Control'
        case 'process_flow_control_documentation_categories':
            return 'Documentation Category'
        case 'process_flow_nodes':
            return 'Process Flow Node'
        case 'process_flow_node_inputs':
            return 'Process Flow Node (Input)'
        case 'process_flow_node_outputs':
            return 'Process Flow Node (Output)'
        case 'process_flow_risks':
            return 'Risk'
        case 'systems':
            return 'System'
        case 'vendors':
            return 'Vendor'
        case 'vendor_products':
            return 'Vendor Product'
    }

    return 'UNKNOWN'
}

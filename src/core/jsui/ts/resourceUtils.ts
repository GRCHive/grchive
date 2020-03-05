import { PageParamsStore } from './pageParams'
import {
    createSingleDbUrl
} from './url'

export interface ResourceHandle {
    displayText: string
    resourceUri: string | null
}

export function standardizeResourceType(typ : string) : string {
    switch (typ) {
        case 'database_connection_info':
            return 'Database Connection Info'
        case 'database_resources':
            return 'Database'
        case 'database_sql_query_requests':
            return 'SQL Query Request'
        case 'document_requests':
            return 'Document Request'
        case 'file_metadata':
            return 'Documentation Metadata'
        case 'file_storage':
            return 'Documentation'
        case 'general_legder_accounts':
            return 'GL Account'
        case 'general_ledger_categories':
            return 'GL Category'
        case 'infrastructure_servers':
            return 'GL Servers'
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

export function getResourceHandle(action : string, typ : string, id : number, extra : any) : Promise<ResourceHandle | null> {
    // We can't use the standard retrieval functions because ideally these should still work in the face
    // of deletion.
    return new Promise<ResourceHandle | null>((resolve, reject) => {
        switch (typ) {
            case 'database_connection_info':
                break
            case 'database_resources':
                break
            case 'database_sql_query_requests':
                break
            case 'document_requests':
                break
            case 'file_metadata':
                break
            case 'file_storage':
                break
            case 'general_legder_accounts':
                break
            case 'general_ledger_categories':
                break
            case 'infrastructure_servers':
                break
            case 'process_flows':
                break
            case 'process_flow_controls':
                break
            case 'process_flow_control_documentation_categories':
                break
            case 'process_flow_nodes':
                break
            case 'process_flow_node_inputs':
                break
            case 'process_flow_node_outputs':
                break
            case 'process_flow_risks':
                break
            case 'systems':
                break
            case 'vendors':
                break
            case 'vendor_products':
                break
            default:
                resolve(null)
                break
        }
    }
}

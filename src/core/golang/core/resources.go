package core

import "errors"

const (
	ResourceDatabase         string = "database_resources"
	ResourceDatabaseConn     string = "database_connection_info"
	ResourceSqlQueryMetadata string = "database_sql_metadata"
	ResourceSqlQuery         string = "database_sql_queries"
	ResourceSqlQueryRequest  string = "database_sql_query_requests"
	ResourceDocMetadata      string = "file_metadata"
	ResourceVendor           string = "vendors"
	ResourceVendorProduct    string = "vendor_products"
	ResourceProcessFlow      string = "process_flows"
	ResourceFlowNode         string = "process_flow_nodes"
	ResourceFlowNodeInput    string = "process_flow_node_inputs"
	ResourceFlowNodeOutput   string = "process_flow_node_outputs"
	ResourceFileStorage      string = "file_storage"
	ResourceGLCat            string = "general_ledger_categories"
	ResourceGLAcc            string = "general_ledger_accounts"
	ResourceServer           string = "infrastructure_servers"
	ResourceSystem           string = "systems"
	ResourceDocRequest       string = "document_requests"
	ResourceDocCat           string = "process_flow_control_documentation_categories"
	ResourceRisk             string = "process_flow_risks"
	ResourceControl          string = "process_flow_controls"
	ResourceUser             string = "users"
)

func GetResourceTypeId(in interface{}) (string, int64, error) {

	switch v := in.(type) {
	case Database:
		return ResourceDatabase, v.Id, nil
	case DatabaseConnection:
		return ResourceDatabaseConn, v.Id, nil
	case DbSqlQueryMetadata:
		return ResourceSqlQueryMetadata, v.Id, nil
	case DbSqlQuery:
		return ResourceSqlQuery, v.Id, nil
	case DbSqlQueryRequest:
		return ResourceSqlQueryRequest, v.Id, nil
	case ControlDocumentationFile:
		return ResourceDocMetadata, v.Id, nil
	case Vendor:
		return ResourceVendor, v.Id, nil
	case VendorProduct:
		return ResourceVendorProduct, v.Id, nil
	case ProcessFlow:
		return ResourceProcessFlow, v.Id, nil
	case ProcessFlowNode:
		return ResourceFlowNode, v.Id, nil
	case FileStorageData:
		return ResourceFileStorage, v.Id, nil
	case GeneralLedgerCategory:
		return ResourceGLCat, v.Id, nil
	case GeneralLedgerAccount:
		return ResourceGLAcc, v.Id, nil
	case Server:
		return ResourceServer, v.Id, nil
	case System:
		return ResourceSystem, v.Id, nil
	case DocumentRequest:
		return ResourceDocRequest, v.Id, nil
	case ControlDocumentationCategory:
		return ResourceDocCat, v.Id, nil
	case Risk:
		return ResourceRisk, v.Id, nil
	case Control:
		return ResourceControl, v.Id, nil
	case User:
		return ResourceUser, v.Id, nil
	}

	return "", 0, errors.New("Unsupported resource.")
}

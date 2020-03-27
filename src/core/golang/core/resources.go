package core

import "errors"

type ResourceHandle struct {
	DisplayText string     `json:"displayText"`
	ResourceUri NullString `json:"resourceUri"`
}

const (
	ResourceIdDatabase         string = "database_resources"
	ResourceIdDatabaseConn     string = "database_connection_info"
	ResourceIdSqlQueryMetadata string = "database_sql_metadata"
	ResourceIdSqlQuery         string = "database_sql_queries"
	ResourceIdSqlQueryRequest  string = "database_sql_query_requests"
	ResourceIdDocMetadata      string = "file_metadata"
	ResourceIdVendor           string = "vendors"
	ResourceIdVendorProduct    string = "vendor_products"
	ResourceIdProcessFlow      string = "process_flows"
	ResourceIdFlowNode         string = "process_flow_nodes"
	ResourceIdFlowNodeInput    string = "process_flow_node_inputs"
	ResourceIdFlowNodeOutput   string = "process_flow_node_outputs"
	ResourceIdFileStorage      string = "file_storage"
	ResourceIdGLCat            string = "general_ledger_categories"
	ResourceIdGLAcc            string = "general_ledger_accounts"
	ResourceIdServer           string = "infrastructure_servers"
	ResourceIdSystem           string = "systems"
	ResourceIdDocRequest       string = "document_requests"
	ResourceIdDocCat           string = "process_flow_control_documentation_categories"
	ResourceIdRisk             string = "process_flow_risks"
	ResourceIdControl          string = "process_flow_controls"
	ResourceIdUser             string = "users"
	ResourceIdClientData       string = "client_data"
)

func GetResourceTypeId(in interface{}) (string, int64, error) {
	if in == nil {
		return "", -1, nil
	}

	switch v := in.(type) {
	case Database:
		return ResourceIdDatabase, v.Id, nil
	case DatabaseConnection:
		return ResourceIdDatabaseConn, v.Id, nil
	case DbSqlQueryMetadata:
		return ResourceIdSqlQueryMetadata, v.Id, nil
	case DbSqlQuery:
		return ResourceIdSqlQuery, v.Id, nil
	case DbSqlQueryRequest:
		return ResourceIdSqlQueryRequest, v.Id, nil
	case ControlDocumentationFile:
		return ResourceIdDocMetadata, v.Id, nil
	case Vendor:
		return ResourceIdVendor, v.Id, nil
	case VendorProduct:
		return ResourceIdVendorProduct, v.Id, nil
	case ProcessFlow:
		return ResourceIdProcessFlow, v.Id, nil
	case ProcessFlowNode:
		return ResourceIdFlowNode, v.Id, nil
	case FileStorageData:
		return ResourceIdFileStorage, v.Id, nil
	case GeneralLedgerCategory:
		return ResourceIdGLCat, v.Id, nil
	case GeneralLedgerAccount:
		return ResourceIdGLAcc, v.Id, nil
	case Server:
		return ResourceIdServer, v.Id, nil
	case System:
		return ResourceIdSystem, v.Id, nil
	case DocumentRequest:
		return ResourceIdDocRequest, v.Id, nil
	case ControlDocumentationCategory:
		return ResourceIdDocCat, v.Id, nil
	case Risk:
		return ResourceIdRisk, v.Id, nil
	case Control:
		return ResourceIdControl, v.Id, nil
	case User:
		return ResourceIdUser, v.Id, nil
	case ClientData:
		return ResourceIdClientData, v.Id, nil
	}

	return "", 0, errors.New("Unsupported resource (GetResourceTypeId).")
}

package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"math"
	"net/http"
	"strconv"
)

type AllAuditTrailInputs struct {
	OrgId int32 `webcore:"orgId"`
}

func allAuditTrailEvents(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllAuditTrailInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	events, err := database.AllFilteredAuditEvents(inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get audit trail events: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(events)
}

type GetAuditTrailInputs struct {
	OrgId              int32          `webcore:"orgId"`
	ResourceHandleOnly bool           `webcore:"resourceHandleOnly"`
	EntryId            core.NullInt64 `webcore:"entryId,optional"`
}

func getAuditTrailEntry(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetAuditTrailInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	org, err := database.FindOrganizationFromId(inputs.OrgId)
	if err != nil {
		core.Warning("Failed to find org: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var event *core.AuditEvent

	if inputs.EntryId.NullInt64.Valid {
		event, err = database.GetAuditEvent(inputs.EntryId.NullInt64.Int64, role)
	} else {
		err = errors.New("Invalid inputs to get audit event.")
	}

	if err != nil {
		core.Warning("Failed to get audit event: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	latestEvent, err := database.GetLatestAuditEvent(event.ResourceType, event.ResourceId, role)
	if err != nil {
		core.Warning("Failed to get latest audit event: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var latestData map[string]interface{}
	latestData, err = database.GetLatestAuditModificationHistoryData(event.ResourceType, event.ResourceId, role)
	if err != nil {
		core.Warning("Failed to get latest modification history: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resourceStillExists := latestEvent.Action != "DELETE"

	if inputs.ResourceHandleOnly {
		type ResourceHandle struct {
			DisplayText string          `json:"displayText"`
			ResourceUri core.NullString `json:"resourceUri"`
		}
		var handle ResourceHandle

		switch event.ResourceType {
		case "database_connection_info":
			dbId := int64(math.Round(event.ResourceExtraData["db_id"].(float64)))
			dbData, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceDatabase,
				dbId,
				role,
			)

			if err != nil {
				core.Warning("Failed to get DB: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if dbData != nil {
				handle.DisplayText = fmt.Sprintf(
					"#%d for DB %s #%d",
					event.ResourceId,
					dbData["name"].(string),
					dbId,
				)

				handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
					webcore.SingleDatabaseRouteName,
					core.DashboardOrgOrgQueryId, org.OktaGroupName,
					core.DashboardOrgDbQueryId, strconv.FormatInt(dbId, 10),
				))
			} else {
				handle.DisplayText = "UNKNOWN"
			}
		case "database_resources":
			handle.DisplayText = fmt.Sprintf(
				"%s #%d",
				latestData["name"].(string),
				event.ResourceId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleDatabaseRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgDbQueryId, strconv.FormatInt(event.ResourceId, 10),
			))
		case "database_sql_query_requests":
			metadataId := int64(math.Round(event.ResourceExtraData["sql_metadata_id"].(float64)))
			metadata, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceSqlQueryMetadata,
				metadataId,
				role,
			)

			if err != nil {
				core.Warning("Failed to get SQL metadata: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			queryId := int64(math.Round(event.ResourceExtraData["query_id"].(float64)))
			query, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceSqlQuery,
				queryId,
				role,
			)

			if err != nil {
				core.Warning("Failed to get SQL query: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if metadata != nil && query != nil {
				handle.DisplayText = fmt.Sprintf(
					"%s #%d for Query %s v%d #%d",
					latestData["name"].(string),
					event.ResourceId,
					metadata["name"].(string),
					int32(math.Round(query["version_number"].(float64))),
					queryId,
				)

				handle.ResourceUri = core.CreateNullString(
					webcore.MustGetRouteUrlAbsolute(
						webcore.SingleDatabaseRouteName,
						core.DashboardOrgOrgQueryId, org.OktaGroupName,
						core.DashboardOrgDbQueryId, strconv.FormatInt(int64(math.Round(metadata["db_id"].(float64))), 10),
					))
			} else {
				handle.DisplayText = "UNKNOWN"
			}
		case "database_sql_metadata":
			dbId := int64(math.Round(event.ResourceExtraData["db_id"].(float64)))
			dbData, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceDatabase,
				dbId,
				role,
			)

			if err != nil {
				core.Warning("Failed to get DB: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if dbData != nil {
				handle.DisplayText = fmt.Sprintf(
					"%s #%d for DB %s #%d",
					latestData["name"].(string),
					event.ResourceId,
					dbData["name"].(string),
					dbId,
				)

				handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
					webcore.SingleDatabaseRouteName,
					core.DashboardOrgOrgQueryId, org.OktaGroupName,
					core.DashboardOrgDbQueryId, strconv.FormatInt(dbId, 10),
				))
			} else {
				handle.DisplayText = "UNKNOWN"
			}
		case "database_sql_queries":
			dbId := int64(math.Round(event.ResourceExtraData["db_id"].(float64)))
			dbData, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceDatabase,
				dbId,
				role,
			)

			if err != nil {
				core.Warning("Failed to get DB: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			metadataId := int64(math.Round(event.ResourceExtraData["sql_metadata_id"].(float64)))
			metadata, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceSqlQueryMetadata,
				metadataId,
				role,
			)

			if err != nil {
				core.Warning("Failed to get metadata: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if dbData != nil {
				handle.DisplayText = fmt.Sprintf(
					"%s v%d #%d for DB %s #%d",
					metadata["name"].(string),
					int32(math.Round(latestData["version_number"].(float64))),
					event.ResourceId,
					dbData["name"].(string),
					dbId,
				)

				handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
					webcore.SingleDatabaseRouteName,
					core.DashboardOrgOrgQueryId, org.OktaGroupName,
					core.DashboardOrgDbQueryId, strconv.FormatInt(dbId, 10),
				))
			} else {
				handle.DisplayText = "UNKNOWN"
			}
		case "document_requests":
			handle.DisplayText = fmt.Sprintf(
				"%s #%d",
				latestData["name"].(string),
				event.ResourceId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleDocRequestRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgDocRequestQueryId, strconv.FormatInt(event.ResourceId, 10),
			))
		case "file_metadata":
			handle.DisplayText = fmt.Sprintf(
				"%s #%d",
				latestData["alt_name"].(string),
				event.ResourceId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleDocumentationRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgDocFileQueryId, strconv.FormatInt(event.ResourceId, 10),
			))
		case "file_storage":
			fileId := int64(math.Round(event.ResourceExtraData["file_id"].(float64)))
			fileData, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceDocMetadata,
				fileId,
				role,
			)

			if err != nil {
				core.Warning("Failed to get file metadata: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			auxData, err := database.GetFileStorageAuxData(event.ResourceId, inputs.OrgId, core.ServerRole)
			if err != nil {
				core.Warning("Failed to get aux file storage data: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			handle.DisplayText = fmt.Sprintf(
				"%s v%d #%d",
				fileData["alt_name"].(string),
				auxData.VersionNumber,
				fileId,
			)

			if auxData.IsPreview {
				handle.DisplayText = handle.DisplayText + " (PREVIEW)"
			}

			handle.ResourceUri = core.CreateNullString(fmt.Sprintf("%s?version=%d", webcore.MustGetRouteUrlAbsolute(
				webcore.SingleDocumentationRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgDocFileQueryId, strconv.FormatInt(event.ResourceId, 10),
			), auxData.VersionNumber))
		case "general_ledger_accounts":
			handle.DisplayText = fmt.Sprintf(
				"%s (%s) #%d",
				latestData["account_name"].(string),
				latestData["account_identifier"].(string),
				event.ResourceId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleGLAccountRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgGLAccQueryId, strconv.FormatInt(event.ResourceId, 10),
			))
		case "general_ledger_categories":
			handle.DisplayText = fmt.Sprintf(
				"%s #%d",
				latestData["name"].(string),
				event.ResourceId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.FullGLAccountRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
			))
		case "infrastructure_servers":
			handle.DisplayText = fmt.Sprintf(
				"%s #%d",
				latestData["name"].(string),
				event.ResourceId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleServerRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgServerQueryId, strconv.FormatInt(event.ResourceId, 10),
			))
		case "process_flows":
			handle.DisplayText = fmt.Sprintf(
				"%s #%d",
				latestData["name"].(string),
				event.ResourceId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleFlowRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgFlowQueryId, strconv.FormatInt(event.ResourceId, 10),
			))
		case "process_flow_controls":
			handle.DisplayText = fmt.Sprintf(
				"%s #%d",
				latestData["name"].(string),
				event.ResourceId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleControlRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgControlQueryId, strconv.FormatInt(event.ResourceId, 10),
			))
		case "process_flow_control_documentation_categories":
			handle.DisplayText = fmt.Sprintf(
				"%s #%d",
				latestData["name"].(string),
				event.ResourceId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleDocCatRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgDocCatQueryId, strconv.FormatInt(event.ResourceId, 10),
			))
		case "process_flow_nodes":
			flowId := int64(math.Round(event.ResourceExtraData["process_flow_id"].(float64)))
			flowData, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceProcessFlow,
				flowId,
				role,
			)

			if err != nil {
				core.Warning("Failed to get flow: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			handle.DisplayText = fmt.Sprintf(
				"%s #%d (%s #%d)",
				latestData["name"].(string),
				event.ResourceId,
				flowData["name"].(string),
				flowId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleFlowRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgFlowQueryId, strconv.FormatInt(flowId, 10),
			))
		case "process_flow_node_inputs":
			flowId := int64(math.Round(event.ResourceExtraData["process_flow_id"].(float64)))
			flowData, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceProcessFlow,
				flowId,
				role,
			)

			if err != nil {
				core.Warning("Failed to get flow: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			nodeId := int64(math.Round(event.ResourceExtraData["node_id"].(float64)))
			nodeData, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceFlowNode,
				nodeId,
				role,
			)

			if err != nil {
				core.Warning("Failed to get node: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			handle.DisplayText = fmt.Sprintf(
				"%s #%d [%s #%d (%s #%d)]",
				latestData["name"].(string),
				event.ResourceId,
				nodeData["name"].(string),
				nodeId,
				flowData["name"].(string),
				flowId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleFlowRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgFlowQueryId, strconv.FormatInt(flowId, 10),
			))
		case "process_flow_node_outputs":
			flowId := int64(math.Round(event.ResourceExtraData["process_flow_id"].(float64)))
			flowData, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceProcessFlow,
				flowId,
				role,
			)

			if err != nil {
				core.Warning("Failed to get flow: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			nodeId := int64(math.Round(event.ResourceExtraData["node_id"].(float64)))
			nodeData, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceFlowNode,
				nodeId,
				role,
			)

			if err != nil {
				core.Warning("Failed to get node: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			handle.DisplayText = fmt.Sprintf(
				"%s #%d [%s #%d (%s #%d)]",
				latestData["name"].(string),
				event.ResourceId,
				nodeData["name"].(string),
				nodeId,
				flowData["name"].(string),
				flowId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleFlowRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgFlowQueryId, strconv.FormatInt(flowId, 10),
			))
		case "process_flow_risks":
			handle.DisplayText = fmt.Sprintf(
				"%s #%d",
				latestData["name"].(string),
				event.ResourceId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleRiskRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgRiskQueryId, strconv.FormatInt(event.ResourceId, 10),
			))
		case "systems":
			handle.DisplayText = fmt.Sprintf(
				"%s #%d",
				latestData["name"].(string),
				event.ResourceId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleSystemRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgSystemQueryId, strconv.FormatInt(event.ResourceId, 10),
			))
		case "vendors":
			handle.DisplayText = fmt.Sprintf(
				"%s #%d",
				latestData["name"].(string),
				event.ResourceId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleVendorRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgVendorQueryId, strconv.FormatInt(event.ResourceId, 10),
			))
		case "vendor_products":
			vendorId := int64(math.Round(event.ResourceExtraData["vendor_id"].(float64)))
			vendorData, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceVendor,
				vendorId,
				role,
			)

			if err != nil {
				core.Warning("Failed to get vendor: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			handle.DisplayText = fmt.Sprintf(
				"%s #%d (%s #%d)",
				latestData["product_name"].(string),
				event.ResourceId,
				vendorData["name"].(string),
				vendorId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleVendorRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgVendorQueryId, strconv.FormatInt(vendorId, 10),
			))

		}

		handle.ResourceUri.NullString.Valid = resourceStillExists
		jsonWriter.Encode(struct {
			Handle ResourceHandle
		}{
			Handle: handle,
		})
	} else {
		core.Warning("Non-resource handle option not yet supported.")
		w.WriteHeader(http.StatusBadRequest)
	}
}

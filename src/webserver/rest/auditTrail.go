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
	OrgId        int32                     `webcore:"orgId"`
	Page         int32                     `webcore:"page"`
	NumItems     int32                     `webcore:"numItems"`
	SortHeader   core.NullString           `webcore:"sortHeaders,optional"`
	SortDesc     core.NullBool             `webcore:"sortDesc,optional"`
	Filter       core.AuditTrailFilterData `webcore:"filter"`
	ResourceType []string                  `webcore:"resourceType,optional"`
	ResourceId   []string                  `webcore:"resourceId,optional"`
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

	sortParams := core.AuditTrailSortParams{
		Limit: core.CreateNullInt32(inputs.NumItems),
		Page:  core.CreateNullInt32(inputs.Page),
	}

	if inputs.SortHeader.NullString.Valid {
		dir := "ASC"
		if inputs.SortDesc.NullBool.Valid && inputs.SortDesc.NullBool.Bool {
			dir = "DESC"
		}
		sortParams.SortDirection = core.CreateNullString(dir)

		switch h := inputs.SortHeader.NullString.String; h {
		case "time":
			sortParams.SortColumns = []string{"hist.performed_at"}
		case "user":
			sortParams.SortColumns = []string{"u.FirstName, u.LastName, u.Email"}
		case "gaction":
			sortParams.SortColumns = []string{"hist.action"}
		case "type":
			sortParams.SortColumns = []string{"hist.resource_type"}
		}
	}

	var events []*core.AuditEvent

	if len(inputs.ResourceType) > 0 && len(inputs.ResourceId) > 0 {
		events, err = database.AllFilteredAuditEventsForResource(
			inputs.ResourceType,
			inputs.ResourceId,
			sortParams,
			inputs.Filter,
			role)
	} else {
		events, err = database.AllFilteredAuditEvents(inputs.OrgId, sortParams, inputs.Filter, role)
	}
	if err != nil {
		core.Warning("Failed to get audit trail events: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var total int

	if len(inputs.ResourceType) > 0 && len(inputs.ResourceId) > 0 {
		total, err = database.CountFilteredAuditEventsForResource(
			inputs.ResourceType,
			inputs.ResourceId,
			inputs.Filter,
			role)
	} else {
		total, err = database.CountFilteredAuditEvents(inputs.OrgId, inputs.Filter, role)
	}

	if err != nil {
		core.Warning("Failed to get total audit events: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(struct {
		Entries []*core.AuditEvent
		Total   int
	}{
		Entries: events,
		Total:   total,
	})
}

type GetAuditTrailInputs struct {
	OrgId   int32          `webcore:"orgId"`
	EntryId core.NullInt64 `webcore:"entryId,optional"`
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

	var handle core.ResourceHandle

	if latestData != nil {
		switch event.ResourceType {
		case "database_connection_info":
			dbId := int64(math.Round(event.ResourceExtraData["db_id"].(float64)))
			dbData, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceIdDatabase,
				strconv.FormatInt(dbId, 10),
				role,
			)

			if err != nil {
				core.Warning("Failed to get DB: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if dbData != nil {
				handle.DisplayText = fmt.Sprintf(
					"#%s for DB %s #%d",
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
				"%s #%s",
				latestData["name"].(string),
				event.ResourceId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleDatabaseRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgDbQueryId, event.ResourceId,
			))
		case "database_sql_query_requests":
			metadataId := int64(math.Round(event.ResourceExtraData["sql_metadata_id"].(float64)))
			metadata, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceIdSqlQueryMetadata,
				strconv.FormatInt(metadataId, 10),
				role,
			)

			if err != nil {
				core.Warning("Failed to get SQL metadata: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			queryId := int64(math.Round(event.ResourceExtraData["query_id"].(float64)))
			query, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceIdSqlQuery,
				strconv.FormatInt(queryId, 10),
				role,
			)

			if err != nil {
				core.Warning("Failed to get SQL query: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if metadata != nil && query != nil {
				handle.DisplayText = fmt.Sprintf(
					"%s #%s for Query %s v%d #%d",
					latestData["name"].(string),
					event.ResourceId,
					metadata["name"].(string),
					int32(math.Round(query["version_number"].(float64))),
					queryId,
				)

				handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
					webcore.SingleSqlRequestRouteName,
					core.DashboardOrgOrgQueryId, org.OktaGroupName,
					core.DashboardOrgSqlRequestQueryId, event.ResourceId,
				))
			} else {
				handle.DisplayText = "UNKNOWN"
			}
		case "database_sql_metadata":
			dbId := int64(math.Round(event.ResourceExtraData["db_id"].(float64)))
			dbData, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceIdDatabase,
				strconv.FormatInt(dbId, 10),
				role,
			)

			if err != nil {
				core.Warning("Failed to get DB: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if dbData != nil {
				handle.DisplayText = fmt.Sprintf(
					"%s #%s for DB %s #%d",
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
				core.ResourceIdDatabase,
				strconv.FormatInt(dbId, 10),
				role,
			)

			if err != nil {
				core.Warning("Failed to get DB: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			metadataId := int64(math.Round(event.ResourceExtraData["sql_metadata_id"].(float64)))
			metadata, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceIdSqlQueryMetadata,
				strconv.FormatInt(metadataId, 10),
				role,
			)

			if err != nil {
				core.Warning("Failed to get metadata: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if dbData != nil && metadata != nil {
				handle.DisplayText = fmt.Sprintf(
					"%s v%d #%s for DB %s #%d",
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
				"%s #%s",
				latestData["name"].(string),
				event.ResourceId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleDocRequestRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgDocRequestQueryId, event.ResourceId,
			))
		case "file_metadata":
			handle.DisplayText = fmt.Sprintf(
				"%s #%s",
				latestData["alt_name"].(string),
				event.ResourceId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleDocumentationRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgDocFileQueryId, event.ResourceId,
			))
		case "file_storage":
			fileId := int64(math.Round(event.ResourceExtraData["file_id"].(float64)))
			fileData, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceIdDocMetadata,
				strconv.FormatInt(fileId, 10),
				role,
			)

			if err != nil {
				core.Warning("Failed to get file metadata: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			storageId, err := strconv.ParseInt(event.ResourceId, 10, 64)
			if err != nil {
				core.Warning("Failed to get storage ID: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			auxData, err := database.GetFileStorageAuxData(storageId, inputs.OrgId, core.ServerRole)
			if fileData != nil && err == nil {
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
					core.DashboardOrgDocFileQueryId, event.ResourceId,
				), auxData.VersionNumber))
			} else {
				handle.DisplayText = "UNKNOWN"
			}
		case "general_ledger_accounts":
			handle.DisplayText = fmt.Sprintf(
				"%s (%s) #%s",
				latestData["account_name"].(string),
				latestData["account_identifier"].(string),
				event.ResourceId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleGLAccountRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgGLAccQueryId, event.ResourceId,
			))
		case "general_ledger_categories":
			handle.DisplayText = fmt.Sprintf(
				"%s #%s",
				latestData["name"].(string),
				event.ResourceId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.FullGLAccountRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
			))
		case "infrastructure_servers":
			handle.DisplayText = fmt.Sprintf(
				"%s #%s",
				latestData["name"].(string),
				event.ResourceId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleServerRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgServerQueryId, event.ResourceId,
			))
		case "process_flows":
			handle.DisplayText = fmt.Sprintf(
				"%s #%s",
				latestData["name"].(string),
				event.ResourceId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleFlowRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgFlowQueryId, event.ResourceId,
			))
		case "process_flow_controls":
			handle.DisplayText = fmt.Sprintf(
				"%s #%s",
				latestData["name"].(string),
				event.ResourceId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleControlRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgControlQueryId, event.ResourceId,
			))
		case "process_flow_control_documentation_categories":
			handle.DisplayText = fmt.Sprintf(
				"%s #%s",
				latestData["name"].(string),
				event.ResourceId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleDocCatRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgDocCatQueryId, event.ResourceId,
			))
		case "process_flow_nodes":
			flowId := int64(math.Round(event.ResourceExtraData["process_flow_id"].(float64)))
			flowData, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceIdProcessFlow,
				strconv.FormatInt(flowId, 10),
				role,
			)

			if err != nil {
				core.Warning("Failed to get flow: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if flowData != nil {
				handle.DisplayText = fmt.Sprintf(
					"%s #%s (%s #%d)",
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
			} else {
				handle.DisplayText = "UNKNOWN"
			}
		case "process_flow_node_inputs":
			flowId := int64(math.Round(event.ResourceExtraData["process_flow_id"].(float64)))
			flowData, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceIdProcessFlow,
				strconv.FormatInt(flowId, 10),
				role,
			)

			if err != nil {
				core.Warning("Failed to get flow: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			nodeId := int64(math.Round(event.ResourceExtraData["node_id"].(float64)))
			nodeData, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceIdFlowNode,
				strconv.FormatInt(nodeId, 10),
				role,
			)

			if err != nil {
				core.Warning("Failed to get node: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if flowData != nil && nodeData != nil {
				handle.DisplayText = fmt.Sprintf(
					"%s #%s [%s #%d (%s #%d)]",
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
			} else {
				handle.DisplayText = "UNKNOWN"
			}
		case "process_flow_node_outputs":
			flowId := int64(math.Round(event.ResourceExtraData["process_flow_id"].(float64)))
			flowData, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceIdProcessFlow,
				strconv.FormatInt(flowId, 10),
				role,
			)

			if err != nil {
				core.Warning("Failed to get flow: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			nodeId := int64(math.Round(event.ResourceExtraData["node_id"].(float64)))
			nodeData, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceIdFlowNode,
				strconv.FormatInt(nodeId, 10),
				role,
			)

			if err != nil {
				core.Warning("Failed to get node: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if flowData != nil && nodeData != nil {
				handle.DisplayText = fmt.Sprintf(
					"%s #%s [%s #%d (%s #%d)]",
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
			} else {
				handle.DisplayText = "UNKNOWN"
			}
		case "process_flow_risks":
			handle.DisplayText = fmt.Sprintf(
				"%s #%s",
				latestData["name"].(string),
				event.ResourceId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleRiskRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgRiskQueryId, event.ResourceId,
			))
		case "systems":
			handle.DisplayText = fmt.Sprintf(
				"%s #%s",
				latestData["name"].(string),
				event.ResourceId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleSystemRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgSystemQueryId, event.ResourceId,
			))
		case "vendors":
			handle.DisplayText = fmt.Sprintf(
				"%s #%s",
				latestData["name"].(string),
				event.ResourceId,
			)

			handle.ResourceUri = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
				webcore.SingleVendorRouteName,
				core.DashboardOrgOrgQueryId, org.OktaGroupName,
				core.DashboardOrgVendorQueryId, event.ResourceId,
			))
		case "vendor_products":
			vendorId := int64(math.Round(event.ResourceExtraData["vendor_id"].(float64)))
			vendorData, err := database.GetLatestAuditModificationHistoryData(
				core.ResourceIdVendor,
				strconv.FormatInt(vendorId, 10),
				role,
			)

			if err != nil {
				core.Warning("Failed to get vendor: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if vendorData != nil {
				handle.DisplayText = fmt.Sprintf(
					"%s #%s (%s #%d)",
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
			} else {
				handle.DisplayText = "UNKNOWN"
			}
		case "client_data":
			dataId, err := strconv.ParseInt(event.ResourceId, 10, 64)
			if err != nil {
				core.Warning("Failed to get client data ID: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			hd, err := webcore.GetResourceHandle(
				event.ResourceType,
				dataId,
				org.Id,
			)

			if err != nil {
				core.Warning("Failed to get client data resource URI: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			handle = *hd
		}
	} else {
		handle.DisplayText = "UNKNOWN"
	}

	var diff map[string]interface{}
	if event.Action == "UPDATE" {
		diff, err = database.GetModificationDiffFromEventId(event.Id, role)
		if err != nil {
			core.Warning("Failed to get audit event diff: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	handle.ResourceUri.NullString.Valid = resourceStillExists
	jsonWriter.Encode(struct {
		Handle core.ResourceHandle
		Diff   map[string]interface{}
	}{
		Handle: handle,
		Diff:   diff,
	})
}

package database

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
)

func FindOrganizationFromId(orgId int32) (*core.Organization, error) {
	rows, err := dbConn.Queryx(`
		SELECT * FROM organizations WHERE id = $1
	`, orgId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var org *core.Organization = new(core.Organization)
	rows.Next()
	err = rows.StructScan(org)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func FindOrganizationFromGroupName(groupName string) (*core.Organization, error) {
	rows, err := dbConn.Queryx(`
		SELECT * FROM organizations WHERE org_group_name = $1
	`, groupName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var org *core.Organization = new(core.Organization)
	rows.Next()
	err = rows.StructScan(org)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func FindOrganizationFromSamlIdP(samlIdp string) (*core.Organization, error) {
	rows, err := dbConn.Queryx(`
		SELECT * FROM organizations WHERE saml_iden = $1
	`, samlIdp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var org *core.Organization = new(core.Organization)
	if !rows.Next() {
		return nil, nil
	}
	err = rows.StructScan(org)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func FindOrganizationFromProcessFlowId(flowId int64, role *core.Role) (*core.Organization, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	org := core.Organization{}
	err := dbConn.Get(&org, `
		SELECT org.*
		FROM process_flows AS pf
		INNER JOIN organizations AS org
			ON org.id = pf.org_id
		WHERE pf.id = $1
	`, flowId)
	return &org, err
}

func FindOrganizationFromNodeId(nodeId int64, role *core.Role) (*core.Organization, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	org := core.Organization{}
	err := dbConn.Get(&org, `
		SELECT org.*
		FROM process_flow_nodes AS node
		INNER JOIN process_flows AS pf
			ON node.process_flow_id = pf.id
		INNER JOIN organizations AS org
			ON org.id = pf.org_id
		WHERE node.id = $1
	`, nodeId)
	return &org, err
}

func FindOrganizationFromRiskId(riskId int64, role *core.Role) (*core.Organization, error) {
	if !role.Permissions.HasAccess(core.ResourceRisks, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	org := core.Organization{}
	err := dbConn.Get(&org, `
		SELECT org.*
		FROM process_flow_risks AS risk
		INNER JOIN organizations AS org
			ON org.id = risk.org_id
		WHERE risk.id = $1
	`, riskId)
	return &org, err
}

func FindOrganizationFromControlId(controlId int64, role *core.Role) (*core.Organization, error) {
	if !role.Permissions.HasAccess(core.ResourceControls, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	org := core.Organization{}
	err := dbConn.Get(&org, `
		SELECT org.*
		FROM process_flow_controls AS ctrl
		INNER JOIN organizations AS org
			ON org.id = ctrl.org_id
		WHERE ctrl.id = $1
	`, controlId)
	return &org, err
}

func FindOrganizationFromDocCatId(catId int64, role *core.Role) (*core.Organization, error) {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentation, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	org := core.Organization{}
	err := dbConn.Get(&org, `
		SELECT org.*
		FROM process_flow_control_documentation_categories AS doc
		INNER JOIN process_flow_controls AS ctrl
			ON doc.control_id = ctrl.id
		INNER JOIN organizations AS org
			ON org.id = ctrl.org_id
		WHERE doc.id = $1
	`, catId)
	return &org, err
}

func FindOrganizationFromProcessFlowInputId(inputId int64, role *core.Role) (*core.Organization, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	org := core.Organization{}
	err := dbConn.Get(&org, `
		SELECT org.*
		FROM process_flow_node_inputs AS io
		INNER JOIN process_flow_nodes AS node
			ON node.id = io.parent_node_id
		INNER JOIN process_flows AS flow
			ON flow.id = node.process_flow_id
		INNER JOIN organizations AS org
			ON org.id = flow.org_id
		WHERE io.id = $1
	`, inputId)
	return &org, err
}

func FindOrganizationFromProcessFlowOutputId(outputId int64, role *core.Role) (*core.Organization, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	org := core.Organization{}
	err := dbConn.Get(&org, `
		SELECT org.*
		FROM process_flow_node_outputs AS io
		INNER JOIN process_flow_nodes AS node
			ON node.id = io.parent_node_id
		INNER JOIN process_flows AS flow
			ON flow.id = node.process_flow_id
		INNER JOIN organizations AS org
			ON org.id = flow.org_id
		WHERE io.id = $1
	`, outputId)
	return &org, err
}

func FindOrganizationFromEdgeId(edgeId int64, role *core.Role) (*core.Organization, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	org := core.Organization{}
	// We assume an organization sanity check was
	// done upon creation of the edge so we can
	// just use the input to grab the organization.
	err := dbConn.Get(&org, `
		SELECT org.*
		FROM process_flow_edges AS edge
		INNER JOIN process_flow_node_inputs AS input
			ON edge.input_id = input.id
		INNER JOIN process_flow_nodes AS node
			ON node.id = input.parent_node_id
		INNER JOIN process_flows AS flow
			ON flow.id = node.process_flow_id
		INNER JOIN organizations AS org
			ON org.id = flow.org_id
		WHERE edge.id = $1
	`, edgeId)
	return &org, err
}

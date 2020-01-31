package core_test

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/grchive/grchive/core"
	"testing"
)

func TestCreateOwnerAccessType(t *testing.T) {
	typ := core.CreateOwnerAccessType()
	assert.NotEqual(t, core.AccessType(0), typ&core.AccessView)
	assert.NotEqual(t, core.AccessType(0), typ&core.AccessEdit)
	assert.NotEqual(t, core.AccessType(0), typ&core.AccessManage)
}

func TestCreateViewOnlyAccessPermission(t *testing.T) {
	p := core.CreateViewOnlyAccessPermission()
	for _, r := range core.AvailableResources {
		typ := p.GetAccessType(r)
		if r == core.ResourceDbConnections {
			// Special case since we don't want to give out passwords for free
			assert.Equal(t, core.AccessType(0), typ&core.AccessView)
		} else {
			assert.NotEqual(t, core.AccessType(0), typ&core.AccessView)
		}
		assert.Equal(t, core.AccessType(0), typ&core.AccessEdit)
		assert.Equal(t, core.AccessType(0), typ&core.AccessManage)
	}
}

func TestCreateAllAccessPermission(t *testing.T) {
	p := core.CreateAllAccessPermission()
	for _, r := range core.AvailableResources {
		typ := p.GetAccessType(r)
		assert.NotEqual(t, core.AccessType(0), typ&core.AccessView)
		assert.NotEqual(t, core.AccessType(0), typ&core.AccessEdit)
		assert.NotEqual(t, core.AccessType(0), typ&core.AccessManage)
	}
}

func TestCreateDefaultRoleMetadata(t *testing.T) {
	m := core.CreateDefaultRoleMetadata(2)
	assert.Equal(t, int32(2), m.OrgId)
	assert.True(t, m.IsDefault)
	assert.False(t, m.IsAdmin)
}

func TestCreateAdminRoleMetadata(t *testing.T) {
	m := core.CreateAdminRoleMetadata(20)
	assert.Equal(t, int32(20), m.OrgId)
	assert.False(t, m.IsDefault)
	assert.True(t, m.IsAdmin)
}

func TestGetAccessType(t *testing.T) {
	p := core.PermissionsMap{}
	for _, r := range core.AvailableResources {
		typ := p.GetAccessType(r)
		assert.Equal(t, core.AccessNone, typ)
	}

	for _, ref := range []struct {
		p core.PermissionsMap
		r core.ResourceType
		a core.AccessType
	}{
		{
			p: core.PermissionsMap{
				OrgUsersAccess: core.AccessManage,
			},
			r: core.ResourceOrgUsers,
			a: core.AccessManage,
		},
		{
			p: core.PermissionsMap{
				OrgRolesAccess: core.AccessEdit,
			},
			r: core.ResourceOrgRoles,
			a: core.AccessEdit,
		},
		{
			p: core.PermissionsMap{
				ProcessFlowsAccess: core.AccessView,
			},
			r: core.ResourceProcessFlows,
			a: core.AccessView,
		},
		{
			p: core.PermissionsMap{
				ControlsAccess: core.AccessView,
			},
			r: core.ResourceControls,
			a: core.AccessView,
		},
		{
			p: core.PermissionsMap{
				ControlDocumentationAccess: core.AccessView,
			},
			r: core.ResourceControlDocumentation,
			a: core.AccessView,
		},
		{
			p: core.PermissionsMap{
				ControlDocMetadataAccess: core.AccessView,
			},
			r: core.ResourceControlDocumentationMetadata,
			a: core.AccessView,
		},
		{
			p: core.PermissionsMap{
				RisksAccess: core.AccessView,
			},
			r: core.ResourceRisks,
			a: core.AccessView,
		},
		{
			p: core.PermissionsMap{
				GLAccess: core.AccessView,
			},
			r: core.ResourceGeneralLedger,
			a: core.AccessView,
		},
		{
			p: core.PermissionsMap{
				SystemAccess: core.AccessView,
			},
			r: core.ResourceSystems,
			a: core.AccessView,
		},
		{
			p: core.PermissionsMap{
				DbAccess: core.AccessView,
			},
			r: core.ResourceDatabases,
			a: core.AccessView,
		},
		{
			p: core.PermissionsMap{
				DbConnectionAccess: core.AccessView,
			},
			r: core.ResourceDbConnections,
			a: core.AccessView,
		},
		{
			p: core.PermissionsMap{
				DocRequestAccess: core.AccessView,
			},
			r: core.ResourceDocRequests,
			a: core.AccessView,
		},
		{
			p: core.PermissionsMap{
				DeploymentAccess: core.AccessView,
			},
			r: core.ResourceDeployments,
			a: core.AccessView,
		},
		{
			p: core.PermissionsMap{
				ServerAccess: core.AccessView,
			},
			r: core.ResourceServers,
			a: core.AccessView,
		},
		{
			p: core.PermissionsMap{
				VendorAccess: core.AccessView,
			},
			r: core.ResourceVendors,
			a: core.AccessView,
		},
	} {
		assert.Equal(t, ref.a, ref.p.GetAccessType(ref.r))
	}
}

func TestHasAccess(t *testing.T) {
	p := core.PermissionsMap{}
	// Can probably for loop this...

	p.OrgUsersAccess = core.AccessNone
	assert.False(t, p.HasAccess(core.ResourceOrgUsers, core.AccessView))
	assert.False(t, p.HasAccess(core.ResourceOrgUsers, core.AccessEdit))
	assert.False(t, p.HasAccess(core.ResourceOrgUsers, core.AccessManage))

	p.OrgUsersAccess = core.AccessView
	assert.True(t, p.HasAccess(core.ResourceOrgUsers, core.AccessView))
	assert.False(t, p.HasAccess(core.ResourceOrgUsers, core.AccessEdit))
	assert.False(t, p.HasAccess(core.ResourceOrgUsers, core.AccessManage))

	p.OrgUsersAccess = core.AccessEdit
	assert.False(t, p.HasAccess(core.ResourceOrgUsers, core.AccessView))
	assert.True(t, p.HasAccess(core.ResourceOrgUsers, core.AccessEdit))
	assert.False(t, p.HasAccess(core.ResourceOrgUsers, core.AccessManage))

	p.OrgUsersAccess = core.AccessManage
	assert.False(t, p.HasAccess(core.ResourceOrgUsers, core.AccessView))
	assert.False(t, p.HasAccess(core.ResourceOrgUsers, core.AccessEdit))
	assert.True(t, p.HasAccess(core.ResourceOrgUsers, core.AccessManage))

	p.OrgUsersAccess = core.AccessView | core.AccessEdit
	assert.True(t, p.HasAccess(core.ResourceOrgUsers, core.AccessView))
	assert.True(t, p.HasAccess(core.ResourceOrgUsers, core.AccessEdit))
	assert.False(t, p.HasAccess(core.ResourceOrgUsers, core.AccessManage))

	p.OrgUsersAccess = core.AccessView | core.AccessManage
	assert.True(t, p.HasAccess(core.ResourceOrgUsers, core.AccessView))
	assert.False(t, p.HasAccess(core.ResourceOrgUsers, core.AccessEdit))
	assert.True(t, p.HasAccess(core.ResourceOrgUsers, core.AccessManage))

	p.OrgUsersAccess = core.AccessEdit | core.AccessManage
	assert.False(t, p.HasAccess(core.ResourceOrgUsers, core.AccessView))
	assert.True(t, p.HasAccess(core.ResourceOrgUsers, core.AccessEdit))
	assert.True(t, p.HasAccess(core.ResourceOrgUsers, core.AccessManage))

	p.OrgUsersAccess = core.AccessView | core.AccessEdit | core.AccessManage
	assert.True(t, p.HasAccess(core.ResourceOrgUsers, core.AccessView))
	assert.True(t, p.HasAccess(core.ResourceOrgUsers, core.AccessEdit))
	assert.True(t, p.HasAccess(core.ResourceOrgUsers, core.AccessManage))
}

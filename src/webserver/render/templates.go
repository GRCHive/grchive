package render

import (
	"encoding/json"
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"html/template"
	"net/http"
)

type templateKey int
type templateMap map[templateKey]*template.Template

var allTemplates templateMap = make(templateMap)

const (
	// Landing Page Keys
	GettingStartedPageTemplateKey templateKey = iota
	ContactUsPageTemplateKey
	LandingPageTemplateKey
	LoginPageTemplateKey
	RegistrationPageTemplateKey
	LearnMorePageTemplateKey
	RedirectTemplateKey
	// Dashboard Keys
	DashboardOrgHomeTemplateKey
	DashboardProcessFlowsTemplateKey
	DashboardSingleFlowTemplateKey
	DashboardUserOrgsTemplateKey
	DashboardUserProfileTemplateKey
	DashboardRisksTemplateKey
	DashboardSingleRiskTemplateKey
	DashboardControlsTemplateKey
	DashboardSingleControlTemplateKey
	DashboardOrgSettingsUsersTemplateKey
	DashboardOrgSettingsRolesTemplateKey
	DashboardOrgSettingsSingleRoleTemplateKey
	DashboardGeneralLedgerTemplateKey
	DashboardGLAccountTemplateKey
	DashboardSystemHomeTemplateKey
	DashboardDbSystemsTemplateKey
	DashboardServersTemplateKey
	DashboardSingleSystemTemplateKey
	DashboardSingleDbTemplateKey
	DashboardSingleServerTemplateKey
	DashboardDocumentationTemplateKey
	DashboardSingleDocumentationTemplateKey
	DashboardDocRequestsTemplateKey
	DashboardSingleDocRequestTemplateKey
	DashboardSingleSqlRequestTemplateKey
	DashboardVendorsTemplateKey
	DashboardSingleVendorTemplateKey
	DashboardSingleDocFileTemplateKey
	DashboardOrgAuditTrailTemplateKey
	// Error Keys
	Error403TemplateKey
	Error404TemplateKey
)

func defaultLoadTemplateWithBase(file string) *template.Template {
	return template.Must(
		template.New("").
			Delims("[[", "]]").
			ParseFiles(
				"src/webserver/templates/base.tmpl",
				"src/core/jsui/main.tmpl",
				file))
}

func defaultLoadTemplateWithDashboardBase(file string) *template.Template {
	return template.Must(
		template.New("").
			Delims("[[", "]]").
			ParseFiles(
				"src/webserver/templates/dashboard/dashboardBase.tmpl",
				"src/core/jsui/dashboard.tmpl",
				file))
}

func RegisterTemplates() {
	// Landing page templates
	allTemplates[GettingStartedPageTemplateKey] =
		defaultLoadTemplateWithBase("src/webserver/templates/gettingStarted.tmpl")
	allTemplates[ContactUsPageTemplateKey] =
		defaultLoadTemplateWithBase("src/webserver/templates/contactUs.tmpl")
	allTemplates[LandingPageTemplateKey] =
		defaultLoadTemplateWithBase("src/webserver/templates/index.tmpl")
	allTemplates[LoginPageTemplateKey] =
		defaultLoadTemplateWithBase("src/webserver/templates/login.tmpl")
	allTemplates[RegistrationPageTemplateKey] =
		defaultLoadTemplateWithBase("src/webserver/templates/register.tmpl")
	allTemplates[LearnMorePageTemplateKey] =
		defaultLoadTemplateWithBase("src/webserver/templates/learnMore.tmpl")
	allTemplates[RedirectTemplateKey] =
		defaultLoadTemplateWithBase("src/webserver/templates/redirect.tmpl")

	// Dashboard templates
	allTemplates[DashboardOrgHomeTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardOrgHome.tmpl")
	allTemplates[DashboardProcessFlowsTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardProcessFlows.tmpl")
	allTemplates[DashboardSingleFlowTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardSingleFlow.tmpl")
	allTemplates[DashboardRisksTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardRisks.tmpl")
	allTemplates[DashboardSingleRiskTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardSingleRisk.tmpl")
	allTemplates[DashboardControlsTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardControls.tmpl")
	allTemplates[DashboardSingleControlTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardSingleControl.tmpl")
	allTemplates[DashboardGeneralLedgerTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardGeneralLedger.tmpl")
	allTemplates[DashboardGLAccountTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardGLAccount.tmpl")
	allTemplates[DashboardSystemHomeTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardSystems.tmpl")
	allTemplates[DashboardDbSystemsTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardDatabases.tmpl")
	allTemplates[DashboardServersTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardServers.tmpl")

	allTemplates[DashboardSingleSystemTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardSingleSystem.tmpl")
	allTemplates[DashboardSingleDbTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardSingleDb.tmpl")
	allTemplates[DashboardSingleServerTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardSingleServer.tmpl")

	allTemplates[DashboardDocumentationTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardDocumentation.tmpl")
	allTemplates[DashboardSingleDocumentationTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardSingleDocumentation.tmpl")
	allTemplates[DashboardSingleDocFileTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardSingleDocFile.tmpl")

	allTemplates[DashboardDocRequestsTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardDocRequests.tmpl")
	allTemplates[DashboardSingleDocRequestTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardSingleDocRequest.tmpl")
	allTemplates[DashboardSingleSqlRequestTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardSingleSqlRequest.tmpl")

	allTemplates[DashboardVendorsTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardVendors.tmpl")
	allTemplates[DashboardSingleVendorTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardSingleVendor.tmpl")

	// Dashboard Org Settings
	allTemplates[DashboardOrgSettingsUsersTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardOrgSettingsUsers.tmpl")
	allTemplates[DashboardOrgSettingsRolesTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardOrgSettingsRoles.tmpl")
	allTemplates[DashboardOrgSettingsSingleRoleTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardOrgSettingsSingleRole.tmpl")

	// Dashboard User
	allTemplates[DashboardUserOrgsTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardUserOrgs.tmpl")
	allTemplates[DashboardUserProfileTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardUserProfile.tmpl")

	// Audit Trail
	allTemplates[DashboardOrgAuditTrailTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardAuditTrail.tmpl")

	// Error templates
	allTemplates[Error403TemplateKey] =
		defaultLoadTemplateWithBase("src/webserver/templates/error/403.tmpl")
	allTemplates[Error404TemplateKey] =
		defaultLoadTemplateWithBase("src/webserver/templates/error/404.tmpl")
}

func RetrieveTemplate(name templateKey) *template.Template {
	if tmpl, ok := allTemplates[name]; ok {
		return tmpl
	}

	core.Error(fmt.Sprintf("Failed to find template: %d", name))
	return nil
}

func RenderTemplate(w http.ResponseWriter, key templateKey, name string, params PageTemplateParameters, extraParams map[string]interface{}) {
	// Handle error?
	jsonRaw, _ := json.Marshal(params)
	RetrieveTemplate(key).
		ExecuteTemplate(
			w,
			name,
			core.MergeMaps(map[string]interface{}{
				"Params": string(jsonRaw),
			}, extraParams))
}

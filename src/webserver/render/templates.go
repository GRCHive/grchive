package render

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"html/template"
	"net/http"
)

type templateKey string
type templateMap map[templateKey]*template.Template

var allTemplates templateMap = make(templateMap)

const (
	// Landing Page Keys
	GettingStartedPageTemplateKey templateKey = "GETSTARTED"
	ContactUsPageTemplateKey                  = "CONTACT"
	LandingPageTemplateKey                    = "LANDING"
	LoginPageTemplateKey                      = "LOGIN"
	RegistrationPageTemplateKey               = "REGISTRATION"
	LearnMorePageTemplateKey                  = "LEARNMORE"
	RedirectTemplateKey                       = "REDIRECT"
	// Dashboard Keys
	DashboardOrgHomeTemplateKey       = "DASHBOARDORGHOME"
	DashboardProcessFlowsTemplateKey  = "DASHBOARDFLOWS"
	DashboardSingleFlowTemplateKey    = "DASHBOARDSINGLEFLOW"
	DashboardUserOrgsTemplateKey      = "DASHBOARDUSERORGS"
	DashboardUserProfileTemplateKey   = "DASHBOARDUSERPROFILE"
	DashboardRisksTemplateKey         = "DASHBOARDRISKS"
	DashboardSingleRiskTemplateKey    = "DASHBOARDSINGLERISK"
	DashboardControlsTemplateKey      = "DASHBOARDCONTROLS"
	DashboardSingleControlTemplateKey = "DASHBOARDSINGLECONTROL"
	// Error Keys
	Error403TemplateKey = "ERROR403"
	Error404TemplateKey = "ERROR404"
)

func defaultLoadTemplateWithBase(file string) *template.Template {
	return template.Must(
		template.New("").
			Delims("[[", "]]").
			ParseFiles("src/webserver/templates/base.tmpl", file))
}

func defaultLoadTemplateWithDashboardBase(file string) *template.Template {
	return template.Must(
		template.New("").
			Delims("[[", "]]").
			ParseFiles("src/webserver/templates/dashboard/dashboardBase.tmpl", file))
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

	// Dashboard User
	allTemplates[DashboardUserOrgsTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardUserOrgs.tmpl")
	allTemplates[DashboardUserProfileTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardUserProfile.tmpl")

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

	core.Error("Failed to find template: " + name)
	return nil
}

func RenderTemplate(w http.ResponseWriter, key templateKey, name string, params PageTemplateParameters) {
	// Handle error?
	jsonRaw, _ := json.Marshal(params)
	RetrieveTemplate(key).
		ExecuteTemplate(
			w,
			name,
			map[string]interface{}{
				"Params": string(jsonRaw),
			})
}

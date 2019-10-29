package render

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
	"html/template"
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
	LearnMorePageTemplateKey                  = "LEARNMORE"
	RedirectTemplateKey                       = "REDIRECT"
	// Dashboard Keys
	DashboardOrgHomeTemplateKey      = "DASHBOARDORGHOME"
	DashboardProcessFlowsTemplateKey = "DASHBOARDFLOWS"
	DashboardUserHomeTemplateKey     = "DASHBOARDUSERHOME"
	DashboardRisksTemplateKey        = "DASHBOARDRISKS"
	DashboardSingleRiskTemplateKey   = "DASHBOARDSINGLERISK"
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
	allTemplates[LearnMorePageTemplateKey] =
		defaultLoadTemplateWithBase("src/webserver/templates/learnMore.tmpl")
	allTemplates[RedirectTemplateKey] =
		defaultLoadTemplateWithBase("src/webserver/templates/redirect.tmpl")

	// Dashboard templates
	allTemplates[DashboardOrgHomeTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardOrgHome.tmpl")
	allTemplates[DashboardUserHomeTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardUserHome.tmpl")
	allTemplates[DashboardProcessFlowsTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardProcessFlows.tmpl")
	allTemplates[DashboardRisksTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardRisks.tmpl")
	allTemplates[DashboardSingleRiskTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardSingleRisk.tmpl")

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

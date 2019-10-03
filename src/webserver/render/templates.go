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
	DashboardHomeTemplateKey = "DASHBOARDHOME"
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

	// Dashing templates
	allTemplates[DashboardHomeTemplateKey] =
		defaultLoadTemplateWithDashboardBase("src/webserver/templates/dashboard/dashboardHome.tmpl")
}

func RetrieveTemplate(name templateKey) *template.Template {
	if tmpl, ok := allTemplates[name]; ok {
		return tmpl
	}

	core.Error("Failed to find template: " + name)
	return nil
}

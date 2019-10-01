package render

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
	"html/template"
)

type templateKey string
type templateMap map[templateKey]*template.Template

var allTemplates templateMap = make(templateMap)

const (
	GettingStartedPageTemplateKey templateKey = "GETSTARTED"
	ContactUsPageTemplateKey      templateKey = "CONTACT"
	LandingPageTemplateKey        templateKey = "LANDING"
	LoginPageTemplateKey          templateKey = "LOGIN"
	LearnMorePageTemplateKey      templateKey = "LEARNMORE"
	GoBackTemplateKey             templateKey = "GOBACK"
)

func defaultLoadTemplateWithBase(file string) *template.Template {
	return template.Must(
		template.New("").
			Delims("[[", "]]").
			ParseFiles("src/webserver/templates/base.tmpl", file))
}

func RegisterTemplates() {
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
	allTemplates[GoBackTemplateKey] =
		defaultLoadTemplateWithBase("src/webserver/templates/goBack.tmpl")
}

func RetrieveTemplate(name templateKey) *template.Template {
	if tmpl, ok := allTemplates[name]; ok {
		return tmpl
	}

	core.Error("Failed to find template: " + name)
	return nil
}

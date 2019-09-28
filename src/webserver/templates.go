package main

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
)

func registerTemplates() {
	allTemplates[GettingStartedPageTemplateKey] =
		template.Must(template.New("").Delims("[[", "]]").ParseFiles("src/webserver/templates/gettingStarted.tmpl"))
	allTemplates[ContactUsPageTemplateKey] =
		template.Must(template.New("").Delims("[[", "]]").ParseFiles("src/webserver/templates/contactUs.tmpl"))
	allTemplates[LandingPageTemplateKey] =
		template.Must(template.New("").Delims("[[", "]]").ParseFiles("src/webserver/templates/index.tmpl"))
	allTemplates[LoginPageTemplateKey] =
		template.Must(template.New("").Delims("[[", "]]").ParseFiles("src/webserver/templates/login.tmpl"))
	allTemplates[LearnMorePageTemplateKey] =
		template.Must(template.New("").Delims("[[", "]]").ParseFiles("src/webserver/templates/learnMore.tmpl"))
}

func retrieveTemplate(name templateKey) *template.Template {
	if tmpl, ok := allTemplates[name]; ok {
		return tmpl
	}

	core.Error("Failed to find template: " + name)
	return nil
}

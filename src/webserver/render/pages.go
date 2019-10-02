package render

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
)

func RenderGettingStartedPage(w http.ResponseWriter, r *http.Request) {
	params, _ := webcore.AddCSRFTokenToRequest(w, r, core.StructToMap(*core.LoadTemplateConfig()))
	RetrieveTemplate(GettingStartedPageTemplateKey).
		ExecuteTemplate(
			w,
			"base",
			params)
}

func RenderContactUsPage(w http.ResponseWriter, r *http.Request) {
	RetrieveTemplate(ContactUsPageTemplateKey).
		ExecuteTemplate(
			w,
			"base",
			core.LoadTemplateConfig())
}

func RenderHomePage(w http.ResponseWriter, r *http.Request) {
	RetrieveTemplate(LandingPageTemplateKey).
		ExecuteTemplate(
			w,
			"base",
			core.LoadTemplateConfig())
}

func RenderLoginPage(w http.ResponseWriter, r *http.Request) {
	params, _ := webcore.AddCSRFTokenToRequest(w, r, core.StructToMap(*core.LoadTemplateConfig()))
	RetrieveTemplate(LoginPageTemplateKey).
		ExecuteTemplate(
			w,
			"base",
			params)
}

func RenderLearnMorePage(w http.ResponseWriter, r *http.Request) {
	RetrieveTemplate(LearnMorePageTemplateKey).
		ExecuteTemplate(
			w,
			"base",
			core.LoadTemplateConfig())
}

func RenderDashboardHomePage(w http.ResponseWriter, r *http.Request) {
}

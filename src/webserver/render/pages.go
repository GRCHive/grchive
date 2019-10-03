package render

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
)

func RenderGettingStartedPage(w http.ResponseWriter, r *http.Request) {
	RetrieveTemplate(GettingStartedPageTemplateKey).
		ExecuteTemplate(
			w,
			"base",
			BuildTemplateParams(w, r, false))
}

func RenderContactUsPage(w http.ResponseWriter, r *http.Request) {
	RetrieveTemplate(ContactUsPageTemplateKey).
		ExecuteTemplate(
			w,
			"base",
			BuildTemplateParams(w, r, false))
}

func RenderHomePage(w http.ResponseWriter, r *http.Request) {
	RetrieveTemplate(LandingPageTemplateKey).
		ExecuteTemplate(
			w,
			"base",
			BuildTemplateParams(w, r, false))
}

func RenderLoginPage(w http.ResponseWriter, r *http.Request) {
	// If the user has a session they can't login...go to dashboard.
	_, err := webcore.FindSessionInContext(r.Context())
	if err == nil {
		http.Redirect(w, r, core.DashboardUrl, http.StatusTemporaryRedirect)
		return
	}

	RetrieveTemplate(LoginPageTemplateKey).
		ExecuteTemplate(
			w,
			"base",
			BuildTemplateParams(w, r, true))
}

func RenderLearnMorePage(w http.ResponseWriter, r *http.Request) {
	RetrieveTemplate(LearnMorePageTemplateKey).
		ExecuteTemplate(
			w,
			"base",
			BuildTemplateParams(w, r, false))
}

func RenderDashboardHomePage(w http.ResponseWriter, r *http.Request) {
	RetrieveTemplate(DashboardHomeTemplateKey).
		ExecuteTemplate(
			w,
			"dashboardBase",
			BuildTemplateParams(w, r, false))
}

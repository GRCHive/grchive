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
		http.Redirect(w, r, webcore.MustGetRouteUrl(webcore.DashboardHomeRouteName), http.StatusTemporaryRedirect)
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
	// For now, redirect user to their organization home page.
	data, err := webcore.FindSessionParsedDataInContext(r.Context())
	if err != nil {
		core.Warning("Bad session parsed data: " + err.Error())
		http.Redirect(w, r,
			webcore.MustGetRouteUrl(webcore.LandingPageRouteName),
			http.StatusTemporaryRedirect)
	} else {
		http.Redirect(w, r,
			webcore.MustGetRouteUrl(
				webcore.DashboardOrgHomeRouteName,
				core.DashboardOrgOrgQueryId,
				data.Org.OktaGroupName),
			http.StatusTemporaryRedirect)
	}
}

func RenderDashboardOrgHomePage(w http.ResponseWriter, r *http.Request) {
	org, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("No organization data: " + err.Error())
		http.Redirect(w, r,
			webcore.MustGetRouteUrl(webcore.DashboardHomeRouteName),
			http.StatusTemporaryRedirect)
		return
	}

	data, err := webcore.FindSessionParsedDataInContext(r.Context())
	if err != nil {
		core.Warning("No user data: " + err.Error())
		http.Redirect(w, r,
			webcore.MustGetRouteUrl(webcore.DashboardHomeRouteName),
			http.StatusTemporaryRedirect)
		return
	}

	RetrieveTemplate(DashboardOrgHomeTemplateKey).
		ExecuteTemplate(
			w,
			"dashboardBase",
			core.MergeMaps(
				BuildTemplateParams(w, r, true),
				BuildOrgTemplateParams(org),
				BuildUserTemplateParams(data.CurrentUser)))
}

func RenderDashboardUserHomePage(w http.ResponseWriter, r *http.Request) {
	data, err := webcore.FindSessionParsedDataInContext(r.Context())
	if err != nil {
		core.Warning("No user data: " + err.Error())
		http.Redirect(w, r,
			webcore.MustGetRouteUrl(webcore.DashboardHomeRouteName),
			http.StatusTemporaryRedirect)
		return
	}

	RetrieveTemplate(DashboardUserHomeTemplateKey).
		ExecuteTemplate(
			w,
			"dashboardBase",
			core.MergeMaps(
				BuildTemplateParams(w, r, true),
				BuildOrgTemplateParams(data.Org),
				BuildUserTemplateParams(data.CurrentUser)))
}

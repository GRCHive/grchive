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
			BuildTemplateParams(w, r))
}

func RenderContactUsPage(w http.ResponseWriter, r *http.Request) {
	RetrieveTemplate(ContactUsPageTemplateKey).
		ExecuteTemplate(
			w,
			"base",
			BuildTemplateParams(w, r))
}

func RenderHomePage(w http.ResponseWriter, r *http.Request) {
	RetrieveTemplate(LandingPageTemplateKey).
		ExecuteTemplate(
			w,
			"base",
			BuildTemplateParams(w, r))
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
			BuildTemplateParams(w, r))
}

func RenderLearnMorePage(w http.ResponseWriter, r *http.Request) {
	RetrieveTemplate(LearnMorePageTemplateKey).
		ExecuteTemplate(
			w,
			"base",
			BuildTemplateParams(w, r))
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

	//data, err := webcore.FindSessionParsedDataInContext(r.Context())
	//if err != nil {
	//	core.Warning("No user data: " + err.Error())
	//	http.Redirect(w, r,
	//		webcore.MustGetRouteUrl(webcore.DashboardHomeRouteName),
	//		http.StatusTemporaryRedirect)
	//	return
	//}

	//RetrieveTemplate(DashboardOrgHomeTemplateKey).
	//	ExecuteTemplate(
	//		w,
	//		"dashboardBase",
	//		core.MergeMaps(
	//			BuildTemplateParams(w, r),
	//			BuildOrgTemplateParams(org),
	//			BuildUserTemplateParams(data.CurrentUser)))

	// TODO: Create some sort of dashboard. For now just redirect to process flows.
	core.Info(webcore.MustGetRouteUrl(webcore.DashboardProcessFlowsRouteName))
	http.Redirect(w, r,
		webcore.MustGetRouteUrl(webcore.DashboardProcessFlowsRouteName, "orgId", org.OktaGroupName),
		http.StatusTemporaryRedirect)
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
				BuildTemplateParams(w, r),
				BuildOrgTemplateParams(data.Org),
				BuildUserTemplateParams(data.CurrentUser)))
}

func RenderDashboardProcessFlowsPage(w http.ResponseWriter, r *http.Request) {
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

	RetrieveTemplate(DashboardProcessFlowsTemplateKey).
		ExecuteTemplate(
			w,
			"dashboardBase",
			core.MergeMaps(
				BuildTemplateParams(w, r),
				BuildOrgTemplateParams(org),
				BuildUserTemplateParams(data.CurrentUser)))
}

func RenderDashboardRisksPage(w http.ResponseWriter, r *http.Request) {
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

	RetrieveTemplate(DashboardRisksTemplateKey).
		ExecuteTemplate(
			w,
			"dashboardBase",
			core.MergeMaps(
				BuildTemplateParams(w, r),
				BuildOrgTemplateParams(org),
				BuildUserTemplateParams(data.CurrentUser)))
}

func RenderDashboardSingleRiskPage(w http.ResponseWriter, r *http.Request) {
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

	RetrieveTemplate(DashboardSingleRiskTemplateKey).
		ExecuteTemplate(
			w,
			"dashboardBase",
			core.MergeMaps(
				BuildTemplateParams(w, r),
				BuildOrgTemplateParams(org),
				BuildUserTemplateParams(data.CurrentUser)))
}

func RenderDashboardControlsPage(w http.ResponseWriter, r *http.Request) {
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

	RetrieveTemplate(DashboardControlsTemplateKey).
		ExecuteTemplate(
			w,
			"dashboardBase",
			core.MergeMaps(
				BuildTemplateParams(w, r),
				BuildOrgTemplateParams(org),
				BuildUserTemplateParams(data.CurrentUser)))
}

func RenderDashboardSingleControlPage(w http.ResponseWriter, r *http.Request) {
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

	RetrieveTemplate(DashboardSingleControlTemplateKey).
		ExecuteTemplate(
			w,
			"dashboardBase",
			core.MergeMaps(
				BuildTemplateParams(w, r),
				BuildOrgTemplateParams(org),
				BuildUserTemplateParams(data.CurrentUser)))
}

func RenderDashboardSingleFlowPage(w http.ResponseWriter, r *http.Request) {
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

	RetrieveTemplate(DashboardSingleFlowTemplateKey).
		ExecuteTemplate(
			w,
			"dashboardBase",
			core.MergeMaps(
				BuildTemplateParams(w, r),
				BuildOrgTemplateParams(org),
				BuildUserTemplateParams(data.CurrentUser)))
}

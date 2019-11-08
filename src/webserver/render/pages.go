package render

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
)

func RenderGettingStartedPage(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, GettingStartedPageTemplateKey, "base", BuildPageTemplateParametersFull(r))
}

func RenderContactUsPage(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, ContactUsPageTemplateKey, "base", BuildPageTemplateParametersFull(r))
}

func RenderHomePage(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, LandingPageTemplateKey, "base", BuildPageTemplateParametersFull(r))
}

func RenderLoginPage(w http.ResponseWriter, r *http.Request) {
	// If the user has a session they can't login...go to dashboard.
	_, err := webcore.FindSessionInContext(r.Context())
	if err == nil {
		http.Redirect(w, r, webcore.MustGetRouteUrl(webcore.DashboardHomeRouteName), http.StatusTemporaryRedirect)
		return
	}

	RenderTemplate(w, LoginPageTemplateKey, "base", BuildPageTemplateParametersFull(r))
}

func RenderLearnMorePage(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, LearnMorePageTemplateKey, "base", BuildPageTemplateParametersFull(r))
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
	http.Redirect(w, r,
		webcore.MustGetRouteUrl(webcore.DashboardProcessFlowsRouteName, "orgId", org.OktaGroupName),
		http.StatusTemporaryRedirect)
}

func verifyContextForOrgDashboard(w http.ResponseWriter, r *http.Request) error {
	_, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("No organization data: " + err.Error())
		http.Redirect(w, r,
			webcore.MustGetRouteUrl(webcore.DashboardHomeRouteName),
			http.StatusTemporaryRedirect)
		return err
	}

	_, err = webcore.FindSessionParsedDataInContext(r.Context())
	if err != nil {
		core.Warning("No user data: " + err.Error())
		http.Redirect(w, r,
			webcore.MustGetRouteUrl(webcore.DashboardHomeRouteName),
			http.StatusTemporaryRedirect)
		return err
	}

	return nil
}

func RenderDashboardUserHomePage(w http.ResponseWriter, r *http.Request) {
	_, err := webcore.FindSessionParsedDataInContext(r.Context())
	if err != nil {
		core.Warning("No user data: " + err.Error())
		http.Redirect(w, r,
			webcore.MustGetRouteUrl(webcore.DashboardHomeRouteName),
			http.StatusTemporaryRedirect)
		return
	}

	RenderTemplate(w, DashboardUserHomeTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r))
}

func RenderDashboardProcessFlowsPage(w http.ResponseWriter, r *http.Request) {
	if verifyContextForOrgDashboard(w, r) != nil {
		return
	}
	RenderTemplate(w, DashboardProcessFlowsTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r))
}

func RenderDashboardRisksPage(w http.ResponseWriter, r *http.Request) {
	if verifyContextForOrgDashboard(w, r) != nil {
		return
	}
	RenderTemplate(w, DashboardRisksTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r))
}

func RenderDashboardSingleRiskPage(w http.ResponseWriter, r *http.Request) {
	if verifyContextForOrgDashboard(w, r) != nil {
		return
	}
	RenderTemplate(w, DashboardSingleRiskTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r))
}

func RenderDashboardControlsPage(w http.ResponseWriter, r *http.Request) {
	if verifyContextForOrgDashboard(w, r) != nil {
		return
	}
	RenderTemplate(w, DashboardControlsTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r))
}

func RenderDashboardSingleControlPage(w http.ResponseWriter, r *http.Request) {
	if verifyContextForOrgDashboard(w, r) != nil {
		return
	}
	RenderTemplate(w, DashboardSingleControlTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r))
}

func RenderDashboardSingleFlowPage(w http.ResponseWriter, r *http.Request) {
	if verifyContextForOrgDashboard(w, r) != nil {
		return
	}
	RenderTemplate(w, DashboardSingleFlowTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r))
}

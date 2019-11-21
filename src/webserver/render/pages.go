package render

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
	"strconv"
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

func RenderRegisterPage(w http.ResponseWriter, r *http.Request) {
	// If the user has a session they can't login...go to dashboard.
	_, err := webcore.FindSessionInContext(r.Context())
	if err == nil {
		http.Redirect(w, r, webcore.MustGetRouteUrl(webcore.DashboardHomeRouteName), http.StatusTemporaryRedirect)
		return
	}

	RenderTemplate(w, RegistrationPageTemplateKey, "base", BuildPageTemplateParametersFull(r))
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
				webcore.DashboardUserHomeRouteName,
				core.DashboardUserQueryId,
				strconv.FormatInt(data.CurrentUser.Id, 10)),
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
	user, err := webcore.FindUserInContext(r.Context())
	if err != nil {
		core.Warning("No user data: " + err.Error())
		http.Redirect(w, r,
			webcore.MustGetRouteUrl(webcore.DashboardHomeRouteName),
			http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(w, r,
		webcore.MustGetRouteUrl(webcore.DashboardUserOrgsRouteName, core.DashboardUserQueryId, strconv.FormatInt(user.Id, 10)),
		http.StatusTemporaryRedirect)
}

func RenderDashboardUserOrgsPage(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardUserOrgsTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r))
}

func RenderDashboardUserProfilePage(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardUserProfileTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r))
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

func RenderDashboardOrgSettingsHome(w http.ResponseWriter, r *http.Request) {
	org, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("No organization data: " + err.Error())
		http.Redirect(w, r,
			webcore.MustGetRouteUrl(webcore.DashboardHomeRouteName),
			http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(w, r,
		webcore.MustGetRouteUrl(webcore.OrgSettingsUsersRouteName, core.DashboardOrgOrgQueryId, org.OktaGroupName),
		http.StatusTemporaryRedirect)
}

func RenderDashboardOrgSettingsUsers(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardOrgSettingsUsersTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r))
}

func RenderDashboardOrgSettingsRoles(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardOrgSettingsRolesTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r))
}

func RenderDashboardOrgSettingsSingleRole(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardOrgSettingsSingleRoleTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r))
}

func RenderDashboardGeneralLedger(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardGeneralLedgerTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r))
}

func RenderDashboardGLAccount(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardGLAccountTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r))
}

func RenderSystemHome(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardSystemHomeTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r))
}

func RenderSingleSystem(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardSingleSystemTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r))
}

func RenderDbSystems(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardDbSystemsTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r))
}

func RenderSingleDb(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardSingleDbTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r))
}

func RenderInfraSystems(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardInfraSystemsTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r))
}

func RenderSingleInfra(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardSingleInfraTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r))
}

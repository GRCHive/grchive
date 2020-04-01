package render

import (
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
	"strconv"
)

var emptyParams = map[string]interface{}{}

func RenderGettingStartedPage(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, GettingStartedPageTemplateKey, "base", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderContactUsPage(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, ContactUsPageTemplateKey, "base", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderHomePage(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, LandingPageTemplateKey, "base", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderLoginPage(w http.ResponseWriter, r *http.Request) {
	// If the user has a session they can't login...go to dashboard.
	_, err := webcore.FindSessionInContext(r.Context())
	if err == nil {
		http.Redirect(w, r, webcore.MustGetRouteUrl(webcore.DashboardHomeRouteName), http.StatusTemporaryRedirect)
		return
	}

	RenderTemplate(w, LoginPageTemplateKey, "base", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderRegisterPage(w http.ResponseWriter, r *http.Request) {
	// If the user has a session they can't login...go to dashboard.
	_, err := webcore.FindSessionInContext(r.Context())
	if err == nil {
		http.Redirect(w, r, webcore.MustGetRouteUrl(webcore.DashboardHomeRouteName), http.StatusTemporaryRedirect)
		return
	}

	RenderTemplate(w, RegistrationPageTemplateKey, "base", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderLearnMorePage(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, LearnMorePageTemplateKey, "base", BuildPageTemplateParametersFull(r), emptyParams)
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
	RenderTemplate(w, DashboardUserOrgsTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderDashboardUserNotificationsPage(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardUserNotificationsTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderDashboardUserProfilePage(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardUserProfileTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderDashboardProcessFlowsPage(w http.ResponseWriter, r *http.Request) {
	if verifyContextForOrgDashboard(w, r) != nil {
		return
	}
	RenderTemplate(w, DashboardProcessFlowsTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderDashboardRisksPage(w http.ResponseWriter, r *http.Request) {
	if verifyContextForOrgDashboard(w, r) != nil {
		return
	}
	RenderTemplate(w, DashboardRisksTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderDashboardSingleRiskPage(w http.ResponseWriter, r *http.Request) {
	if verifyContextForOrgDashboard(w, r) != nil {
		return
	}
	RenderTemplate(w, DashboardSingleRiskTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderDashboardControlsPage(w http.ResponseWriter, r *http.Request) {
	if verifyContextForOrgDashboard(w, r) != nil {
		return
	}
	RenderTemplate(w, DashboardControlsTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderDashboardSingleControlPage(w http.ResponseWriter, r *http.Request) {
	if verifyContextForOrgDashboard(w, r) != nil {
		return
	}
	RenderTemplate(w, DashboardSingleControlTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderDashboardSingleFlowPage(w http.ResponseWriter, r *http.Request) {
	if verifyContextForOrgDashboard(w, r) != nil {
		return
	}
	RenderTemplate(w, DashboardSingleFlowTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
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
	RenderTemplate(w, DashboardOrgSettingsUsersTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderDashboardOrgSettingsRoles(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardOrgSettingsRolesTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderDashboardOrgSettingsSingleRole(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardOrgSettingsSingleRoleTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderDashboardGeneralLedger(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardGeneralLedgerTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderDashboardGLAccount(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardGLAccountTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderSystemHome(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardSystemHomeTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderSingleSystem(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardSingleSystemTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderDbSystems(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardDbSystemsTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderSingleDb(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardSingleDbTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderServers(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardServersTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderSingleServer(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardSingleServerTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderDocumentation(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardDocumentationTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderSingleDocCat(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardSingleDocumentationTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderDocRequest(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardDocRequestsTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderSingleDocRequest(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardSingleDocRequestTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderSingleSqlRequest(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardSingleSqlRequestTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderVendors(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardVendorsTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderSingleVendor(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardSingleVendorTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderSingleDocFile(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardSingleDocFileTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderAuditTrail(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardOrgAuditTrailTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderClientData(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardOrgClientDataTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderSingleClientData(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, DashboardOrgSingleClientDataTemplateKey, "dashboardBase", BuildPageTemplateParametersFull(r), emptyParams)
}

func RenderRedirectPage(w http.ResponseWriter, r *http.Request, url string) {
	RenderTemplate(w, RedirectTemplateKey, "base",
		BuildPageTemplateParametersFull(r),
		CreateRedirectParams(w, r, "Oops!",
			"Something went wrong! Please try again.",
			url))
}

func RenderFeatureRequestPage(w http.ResponseWriter, r *http.Request, featureId core.FeatureId, pending bool) {
	featureName, err := database.GetFeatureName(featureId)
	if err != nil {
		Render404(w, r)
	}

	RenderTemplate(
		w,
		DashboardOrgFeatureRequestTemplateKey,
		"dashboardBase",
		BuildPageTemplateParametersFull(r),
		map[string]interface{}{
			"FeatureName": featureName,
			"FeatureId":   featureId,
			"Pending":     pending,
		},
	)
}

package core

import (
	"fmt"
	"net/url"
)

// Landing Page
const GetStartedUrl string = "/getting-started"
const ContactUsUrl string = "/contact-us"
const HomePageUrl string = "/"
const LoginUrl string = "/login"
const RegisterUrl string = "/register"
const LogoutUrl string = "/logout"
const LearnMoreUrl string = "/learn"

func CreateSamlCallbackUrl() string {
	return EnvConfig.Login.RedirectUrl
}

const VerifyEmailUrl string = "/verify"

// Dashboard
const DashboardUrl string = "/dashboard"
const DashboardHomeUrl string = "/"

// Dashboard - Organization
const DashboardOrgOrgQueryId string = "orgId"
const DashboardOrgRiskQueryId string = "risk"
const DashboardOrgControlQueryId string = "control"
const DashboardOrgFlowQueryId string = "flow"

var DashboardOrgUrl string = fmt.Sprintf("/org/{%s}", DashboardOrgOrgQueryId)

const DashboardOrgHomeUrl string = "/"

const DashboardOrgAllFlowsEndpoint string = "/flows"
const DashboardOrgAllRiskEndpoint string = "/risks"
const DashboardOrgAllControlsEndpoint string = "/controls"

var DashboardOrgRiskEndpoint string = fmt.Sprintf("/risks/{%s}", DashboardOrgRiskQueryId)
var DashboardOrgControlEndpoint string = fmt.Sprintf("/controls/{%s}", DashboardOrgControlQueryId)
var DashboardOrgFlowEndpoint string = fmt.Sprintf("/flows/{%s}", DashboardOrgFlowQueryId)

// Dashboard - User
const DashboardUserQueryId string = "user"

var DashboardUserPrefix string = fmt.Sprintf("/user/{%s}", DashboardUserQueryId)

const DashboardUserHomeUrl string = "/"
const DashboardUserOrgUrl string = "/orgs"
const DashboardUserProfileUrl string = "/profile"

// API
const ApiUrl string = "/api"

// API Verification
const ApiVerificationPrefix string = "/verification"
const ApiResendVerificationEndpoint string = "/resend"

// API - Process Flow
const ProcessFlowQueryId string = "flow"
const ApiProcessFlowUrl string = "/flows"
const ApiProcessFlowGetAllUrl string = "/"
const ApiProcessFlowNewUrl string = "/new"
const ApiProcessFlowDeleteEndpoint string = "/delete"

var ApiProcessFlowUpdateUrl string = fmt.Sprintf("/{%s}/update", ProcessFlowQueryId)
var ApiProcessFlowGetFullDataUrl string = fmt.Sprintf("/{%s}/full", ProcessFlowQueryId)

// API - Process Flow Nodes
const ApiProcessFlowNodesUrl string = "/flownodes"
const ApiProcessFlowNodesGetTypesUrl string = "/types"
const ApiProcessFlowNodesNewUrl string = "/new"
const ApiProcessFlowNodesEditUrl string = "/edit"
const ApiProcessFlowNodesDeleteUrl string = "/delete"

// API - Process Flow Edges
const ApiProcessFlowEdgesUrl string = "/flowedges"
const ApiProcessFlowEdgesNewUrl string = "/new"
const ApiProcessFlowEdgesDeleteUrl string = "/delete"

// API - Process Flow IO
const ApiProcessFlowIOUrl string = "/flowio"
const ApiProcessFlowIOGetTypesUrl string = "/types"
const ApiProcessFlowIONewUrl string = "/new"
const ApiProcessFlowIODeleteUrl string = "/delete"
const ApiProcessFlowIOEditUrl string = "/edit"

// API - Users
var ApiUserUrl string = fmt.Sprintf("/user/{%s}", DashboardUserQueryId)

const ApiUserProfileUrl string = "/profile"

// API - Organization
var ApiOrgPrefix string = fmt.Sprintf("/org/{%s}", DashboardOrgOrgQueryId)

const ApiOrgUsersEndpoint string = "/users"

// API - Risks
const ApiRiskPrefix string = "/risk"
const ApiNewRiskEndpoint string = "/new"
const ApiDeleteRiskEndpoint string = "/delete"
const ApiEditRiskEndpoint string = "/edit"
const ApiAddRiskToNodeEndpoint string = "/add"
const ApiGetAllRisksEndpoint string = "/"

var ApiGetSingleRiskEndpoint string = fmt.Sprintf("/{%s}", DashboardOrgRiskQueryId)

// API - Controls
const ApiControlPrefix string = "/control"
const ApiNewControlEndpoint string = "/new"
const ApiGetControlTypesEndpoint string = "/types"
const ApiDeleteControlEndpoint string = "/delete"
const ApiEditControlEndpoint string = "/edit"
const ApiAddControlEndpoint string = "/add"
const ApiGetAllControlEndpoint string = "/"

var ApiGetSingleControlEndpoint string = fmt.Sprintf("/{%s}", DashboardOrgControlQueryId)

// API - Control Documentation
const ApiControlDocumentationPrefix string = "/documentation"
const ApiNewControlDocumentationCategoryEndpoint string = "/newcat"
const ApiEditControlDocumentationCategoryEndpoint string = "/editcat"
const ApiDeleteControlDocumentationCategoryEndpoint string = "/deletecat"
const ApiUploadControlDocumentationEndpoint string = "/upload"
const ApiGetControlDocumentationEndpoint string = "/get"
const ApiDeleteControlDocumentationEndpoint string = "/delete"
const ApiDownloadControlDocumentationEndpoint string = "/download"

// Websocket
const WebsocketPrefix string = "/ws"

var WebsocketProcessFlowNodeDisplaySettingsEndpoint = fmt.Sprintf("/flownodedisp/{%s}", ProcessFlowQueryId)

func CreateUrlWithParams(input string, params map[string]string) (string, error) {
	gUrl, err := url.Parse(input)
	if err != nil {
		return "", err
	}

	q := gUrl.Query()
	for k, v := range params {
		q.Add(k, v)
	}

	gUrl.RawQuery = q.Encode()
	return gUrl.String(), nil
}

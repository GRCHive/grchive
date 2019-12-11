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
const AcceptInviteUrl string = "/invite"

// Dashboard
const DashboardUrl string = "/dashboard"
const DashboardHomeUrl string = "/"

// Dashboard - Organization
const DashboardOrgOrgQueryId string = "orgId"
const DashboardOrgRiskQueryId string = "risk"
const DashboardOrgControlQueryId string = "control"
const DashboardOrgFlowQueryId string = "flow"
const DashboardOrgRoleQueryId string = "roleId"
const DashboardOrgGLAccQueryId string = "accId"
const DashboardOrgSystemQueryId string = "sysId"
const DashboardOrgDbQueryId string = "dbId"
const DashboardOrgInfraQueryId string = "infraId"
const DashboardOrgDocCatQueryId string = "docCatId"
const DashboardOrgDocRequestQueryId string = "reqId"

var DashboardOrgUrl string = fmt.Sprintf("/org/{%s}", DashboardOrgOrgQueryId)

const DashboardOrgHomeUrl string = "/"

const DashboardOrgAllFlowsEndpoint string = "/flows"
const DashboardOrgAllRiskEndpoint string = "/risks"
const DashboardOrgAllControlsEndpoint string = "/controls"
const DashboardOrgAllDocumentationEndpoint string = "/documentation"
const DashboardOrgAllDocRequestsEndpoint string = "/requests"

var DashboardOrgRiskEndpoint string = fmt.Sprintf("/risks/{%s}", DashboardOrgRiskQueryId)
var DashboardOrgControlEndpoint string = fmt.Sprintf("/controls/{%s}", DashboardOrgControlQueryId)
var DashboardOrgFlowEndpoint string = fmt.Sprintf("/flows/{%s}", DashboardOrgFlowQueryId)
var DashboardOrgSingleDocCatEndpoint string = fmt.Sprintf("/documentation/cat/{%s}", DashboardOrgDocCatQueryId)
var DashboardOrgSingleDocRequestEndpoint string = fmt.Sprintf("/requests/{%s}", DashboardOrgDocRequestQueryId)

// Dashboard - Organization - Settings
const DashboardOrgSettingsPrefix string = "/settings"
const DashboardOrgSettingsHomeEndpoint string = "/"
const DashboardOrgSettingsUsersEndpoint string = "/users"
const DashboardOrgSettingsRolesEndpoint string = "/roles"

var DashboardOrgSettingsSingleRoleEndpoint string = fmt.Sprintf("/roles/{%s}", DashboardOrgRoleQueryId)

// Dashboard - User
const DashboardUserQueryId string = "user"

var DashboardUserPrefix string = fmt.Sprintf("/user/{%s}", DashboardUserQueryId)

const DashboardUserHomeUrl string = "/"
const DashboardUserOrgUrl string = "/orgs"
const DashboardUserProfileUrl string = "/profile"

// Dashboard - General Ledger
const DashboardGeneralLedgerPrefix string = "/gl"
const DashboardGeneralLedgerViewEndpoint string = "/"

var DashboardOrgGLAccountEndpoint string = fmt.Sprintf("/acc/{%s}", DashboardOrgGLAccQueryId)

// Dashboard - Systems
const DashboardSystemsPrefix string = "/it"
const DashboardSystemHomeEndpoint string = "/systems"
const DashboardDbSystemsEndpoint string = "/databases"
const DashboardInfraSystemsEndpoint string = "/infrastructure"

var DashboardSingleSystemEndpoint string = fmt.Sprintf("%s/{%s}", DashboardSystemHomeEndpoint, DashboardOrgSystemQueryId)
var DashboardSingleDbEndpoint string = fmt.Sprintf("%s/{%s}", DashboardDbSystemsEndpoint, DashboardOrgDbQueryId)
var DashboardSingleInfraEndpoint string = fmt.Sprintf("%s/{%s}", DashboardInfraSystemsEndpoint, DashboardOrgInfraQueryId)

// Generic API Actions
const ApiNewEndpoint = "/new"
const ApiAllEndpoint = "/all"
const ApiGetEndpoint = "/get"
const ApiDeleteEndpoint = "/delete"

// API
const ApiUrl string = "/api"

// API - Invites
const ApiInvitePrefix string = "/invite"
const ApiSendInviteEndpoint string = "/send"

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
var ApiUserPrefix string = fmt.Sprintf("/user/{%s}", DashboardUserQueryId)

const ApiUserProfileEndpoint string = "/profile"
const ApiUserOrgsEndpoint string = "/orgs"

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
const ApiLinkDocCatControlEndpoint string = "/linkCat"
const ApiUnlinkDocCatControlEndpoint string = "/unlinkCat"
const ApiGetAllControlEndpoint string = "/"

var ApiGetSingleControlEndpoint string = fmt.Sprintf("/{%s}", DashboardOrgControlQueryId)

// API - Control Documentation
const ApiControlDocumentationPrefix string = "/documentation"
const ApiNewControlDocumentationCategoryEndpoint string = "/cat/new"
const ApiEditControlDocumentationCategoryEndpoint string = "/cat/edit"
const ApiDeleteControlDocumentationCategoryEndpoint string = "/cat/delete"
const ApiAllControlDocumentationCategoryEndpoint string = "/cat/all"
const ApiGetControlDocumentationCategoryEndpoint string = "/cat/get"

const ApiUploadControlDocumentationEndpoint string = "/file/upload"
const ApiGetControlDocumentationEndpoint string = "/file/get"
const ApiDeleteControlDocumentationEndpoint string = "/file/delete"
const ApiDownloadControlDocumentationEndpoint string = "/file/download"

// API - Roles
const ApiRolePrefix string = "/roles"
const ApiGetOrganizationRolesEndpoint string = "/all"
const ApiGetSingleRoleEndpoint string = "/get"
const ApiNewRoleEndpoint string = "/new"
const ApiEditRoleEndpoint string = "/edit"
const ApiDeleteRoleEndpoint string = "/delete"
const ApiAddUsersToRoleEndpoint string = "/addUsers"

// API - General Ledger
const ApiGeneralLedgerPrefix = "/gl"
const ApiGetGLLevelEndpoint string = "/get"
const ApiNewGLCategoryEndpoint string = "/cat/new"
const ApiEditGLCategoryEndpoint string = "/cat/edit"
const ApiDeleteGLCategoryEndpoint string = "/cat/delete"
const ApiNewGLAccountEndpoint string = "/acc/new"
const ApiGetGLAccountEndpoint string = "/acc/get"
const ApiEditGLAccountEndpoint string = "/acc/edit"
const ApiDeleteGLAccountEndpoint string = "/acc/delete"

// API - IT
const ApiITPrefix = "/it"
const ApiITDeleteDbSysLinkEndpoint = "/deleteDbSysLink"

// API - IT - Systems
const ApiITSystemsPrefix = "/systems"
const ApiITSystemsNewEndpoint = "/new"
const ApiITSystemsAllEndpoint = "/all"
const ApiITSystemGetEndpoint = "/get"
const ApiITSystemEditEndpoint = "/edit"
const ApiITSystemDeleteEndpoint = "/delete"
const ApiITSystemLinkDbEndpoint = "/linkDb"

// API - IT - DB
const ApiITDbPrefix = "/db"
const ApiITDbNewEndpoint = "/new"
const ApiITDbAllEndpoint = "/all"
const ApiITDbTypesEndpoint = "/types"
const ApiITDbGetEndpoint = "/get"
const ApiITDbEditEndpoint = "/edit"
const ApiITDbDeleteEndpoint = "/delete"
const ApiITDbLinkSysEndpoint = "/linkSys"

// API - IT - DB CONN
const ApiITDbConnPrefix = "/connection"
const ApiITDbConnNewEndpoint = "/new"
const ApiITDbConnDeleteEndpoint = "/delete"

// API - IT - Infrastructure
const ApiITInfraPrefix = "/infra"

// API - Document Requests
const ApiDocRequestPrefix = "/requests"
const ApiDocRequestCompleteEndpoint = "/complete"

// API - Document Request Comments
const ApiDocRequestCommentsPrefix = "/comments"

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

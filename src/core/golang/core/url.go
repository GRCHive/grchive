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
const DashboardOrgServerQueryId string = "serverId"
const DashboardOrgDocCatQueryId string = "docCatId"
const DashboardOrgDocFileQueryId string = "docFileId"
const DashboardOrgDocRequestQueryId string = "reqId"
const DashboardOrgSqlRequestQueryId string = "sqlReqId"
const DashboardOrgVendorQueryId string = "vendorId"
const DashboardOrgClientDataQueryId string = "clientDataId"
const DashboardOrgClientScriptQueryId string = "clientScriptId"
const DashboardOrgCommitQueryId string = "commitId"
const DashboardOrgScriptRunQueryId string = "runId"
const DashboardOrgScriptRequestQueryId = "scriptReqId"
const DashboardOrgGenericRequestQueryId = "reqId"
const DashboardOrgShellScriptQueryId = "shellScriptId"
const DashboardOrgShellScriptVersionQueryId = "shellScriptVersionId"
const DashboardOrgServerSshConnectionQueryId = "serverSshConnectionId"

var DashboardOrgUrl string = fmt.Sprintf("/org/{%s}", DashboardOrgOrgQueryId)

const DashboardOrgHomeUrl string = "/"

const DashboardOrgAllFlowsEndpoint string = "/flows"
const DashboardOrgAllRiskEndpoint string = "/risks"
const DashboardOrgAllControlsEndpoint string = "/controls"
const DashboardOrgAllDocumentationEndpoint string = "/documentation"
const DashboardOrgAllDocRequestsEndpoint string = "/requests"
const DashboardOrgAllVendorsEndpoint string = "/vendors"

var DashboardOrgRiskEndpoint string = fmt.Sprintf("/risks/{%s}", DashboardOrgRiskQueryId)
var DashboardOrgControlEndpoint string = fmt.Sprintf("/controls/{%s}", DashboardOrgControlQueryId)
var DashboardOrgFlowEndpoint string = fmt.Sprintf("/flows/{%s}", DashboardOrgFlowQueryId)
var DashboardOrgSingleDocCatEndpoint string = fmt.Sprintf("/documentation/cat/{%s}", DashboardOrgDocCatQueryId)
var DashboardOrgSingleDocFileEndpoint string = fmt.Sprintf("/documentation/file/{%s}", DashboardOrgDocFileQueryId)
var DashboardOrgSingleDocRequestEndpoint string = fmt.Sprintf("/requests/doc/{%s}", DashboardOrgDocRequestQueryId)
var DashboardOrgSingleSqlRequestEndpoint string = fmt.Sprintf("/requests/sql/{%s}", DashboardOrgSqlRequestQueryId)
var DashboardOrgSingleScriptRequestEndpoint string = fmt.Sprintf("/requests/scripts/{%s}", DashboardOrgScriptRequestQueryId)
var DashboardOrgSingleVendorEndpoint string = fmt.Sprintf("/vendors/{%s}", DashboardOrgVendorQueryId)
var DashboardSingleShellEndpoint string = fmt.Sprintf("/shell/{%s}", DashboardOrgShellScriptQueryId)

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
const DashboardUserNotificationsUrl string = "/notifications"

// Dashboard - General Ledger
const DashboardGeneralLedgerPrefix string = "/gl"
const DashboardGeneralLedgerViewEndpoint string = "/"

var DashboardOrgGLAccountEndpoint string = fmt.Sprintf("/acc/{%s}", DashboardOrgGLAccQueryId)

// Dashboard - Systems
const DashboardSystemsPrefix string = "/it"
const DashboardSystemHomeEndpoint string = "/systems"
const DashboardDbSystemsEndpoint string = "/databases"
const DashboardServersEndpoint string = "/servers"
const DashboardShellEndpoint string = "/shell"

var DashboardSingleSystemEndpoint string = fmt.Sprintf("%s/{%s}", DashboardSystemHomeEndpoint, DashboardOrgSystemQueryId)
var DashboardSingleDbEndpoint string = fmt.Sprintf("%s/{%s}", DashboardDbSystemsEndpoint, DashboardOrgDbQueryId)
var DashboardSingleServerEndpoint string = fmt.Sprintf("%s/{%s}", DashboardServersEndpoint, DashboardOrgServerQueryId)

// Dashboard - Automation
const DashboardAutomationPrefix string = "/auto"
const DashboardDataPrefix string = "/data"
const DashboardDataSourcePrefix string = "/source"
const DashboardCodePrefix string = "/code"
const DashboardScriptPrefix string = "/scripts"
const DashboardLogsPrefix string = "/logs"
const DashboardStatusPrefix string = "/status"
const DashboardRunPrefix string = "/runs"
const DashboardLinkPrefix string = "/link"

var DashboardSingleDataEndpoint string = fmt.Sprintf("/{%s}", DashboardOrgClientDataQueryId)
var DashboardSingleScriptEndpoint string = fmt.Sprintf("/{%s}", DashboardOrgClientScriptQueryId)

var DashboardSingleBuildLogEndpoint string = fmt.Sprintf("/build/{%s}", DashboardOrgCommitQueryId)
var DashboardSingleScriptRunLogEndpoint string = fmt.Sprintf("/run/{%s}", DashboardOrgScriptRunQueryId)

// Generic API Actions
const ApiNewEndpoint = "/new"
const ApiAllEndpoint = "/all"
const ApiGetEndpoint = "/get"
const ApiDeleteEndpoint = "/delete"
const ApiUpdateEndpoint = "/update"
const ApiDuplicateEndpoint = "/duplicate"
const ApiRunEndpoint = "/run"
const ApiReadEndpoint = "/read"
const ApiSaveEndpoint = "/save"
const ApiLinkEndpoint = "/link"

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

const ApiProcessFlowNodeLinksPrefix = "/link"

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
const ApiProcessFlowIOOrderEndpoint string = "/order"

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
const ApiGetSingleRiskEndpoint string = "/get"

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
const ApiGetSingleControlEndpoint string = "/get"

// API - Control Documentation
const ApiControlDocumentationPrefix string = "/documentation"
const ApiDocCatPrefix = "/cat"
const ApiFolderPrefix = "/folder"
const ApiDocFilePrefix = "/file"

const ApiNewControlDocumentationCategoryEndpoint string = "/cat/new"
const ApiEditControlDocumentationCategoryEndpoint string = "/cat/edit"
const ApiDeleteControlDocumentationCategoryEndpoint string = "/cat/delete"
const ApiAllControlDocumentationCategoryEndpoint string = "/cat/all"
const ApiGetControlDocumentationCategoryEndpoint string = "/cat/get"

const ApiUploadControlDocumentationEndpoint string = "/file/upload"
const ApiAllControlDocumentationEndpoint string = "/file/all"
const ApiDeleteControlDocumentationEndpoint string = "/file/delete"
const ApiDownloadControlDocumentationEndpoint string = "/file/download"
const ApiGetControlDocumentationEndpoint string = "/file/get"
const ApiEditControlDocumentationEndpoint string = "/file/edit"
const ApiRegenPreviewControlDocumentationEndpoint string = "/file/preview"

const ApiFileVersionPrefix string = "/file/versions"

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

// API - IT - DB SQL
const ApiITSqlPrefix = "/sql"
const ApiITSqlRefreshPrefix = "/refresh"
const ApiITSqlSchemaPrefix = "/schema"
const ApiITSqlQueriesPrefix = "/query"
const ApiITSqlRequestsPrefix = "/requests"
const ApiITSqlRequestsStatusEndpoint = "/status"

// API - IT - Server
const ApiITServerPrefix = "/servers"

// API - Document Requests
const ApiDocRequestPrefix = "/requests"
const ApiDocRequestCompleteEndpoint = "/complete"

// API - Comments
const ApiCommentsPrefix = "/comments"

// API - Deployment
const ApiDeploymentPrefix = "/deployment"
const ApiLinkPrefix = "/link"

// API - Vendor
const ApiVendorPrefix = "/vendor"
const ApiVendorProductPrefix = "/product"
const ApiVendorProductSocPrefix = "/soc"

// API - Audit Trail
const ApiAuditTrailPrefix = "/auditlog"

// API - Notification
const ApiNotificationPrefix = "/notifications"

// API - Resource
const ApiResourcePrefix = "/resource"

// API - Features
const ApiFeaturePrefix = "/feature"

// Websocket
const WebsocketPrefix string = "/ws"

// Metadata
const ApiMetadataPrefix = "/metadata"

var WebsocketProcessFlowNodeDisplaySettingsEndpoint = fmt.Sprintf("/flownodedisp/{%s}", ProcessFlowQueryId)
var WebsocketUserNotificationsEndpoint = fmt.Sprintf("/notifications/{%s}", DashboardUserQueryId)

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

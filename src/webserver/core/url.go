package core

import "fmt"

// Landing Page
const GetStartedUrl string = "/getting-started"
const ContactUsUrl string = "/contact-us"
const HomePageUrl string = "/"
const LoginUrl string = "/login"
const LogoutUrl string = "/logout"
const LearnMoreUrl string = "/learn"

var SamlCallbackUrl string = LoadEnvConfig().Login.RedirectUrl

// Dashboard
const DashboardUrl string = "/dashboard"
const DashboardHomeUrl string = "/"

// Dashboard - Organization
const DashboardOrgOrgQueryId string = "orgId"

var DashboardOrgUrl string = fmt.Sprintf("/org/{%s}", DashboardOrgOrgQueryId)

const DashboardOrgHomeUrl string = "/"
const DashboardOrgFlowUrl string = "/flows"

// Dashboard - User
const DashboardUserQueryId string = "user"

var DashboardUserUrl string = fmt.Sprintf("/user/{%s}", DashboardUserQueryId)

const DashboardUserHomeUrl string = "/"

// API
const ApiUrl string = "/api"

// API - Process Flow
const ProcessFlowQueryId string = "flow"
const ApiProcessFlowUrl string = "/flows"
const ApiProcessFlowGetAllUrl string = "/"
const ApiProcessFlowNewUrl string = "/new"

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

// Websocket
const WebsocketPrefix string = "/ws"
const WebsocketProcessFlowNodeDisplaySettingsEndpoint = "/flownodedisp"

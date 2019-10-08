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
const ApiProcessFlowUrl string = "/flows"
const ApiProcessFlowGetAllUrl string = "/"
const ApiProcessFlowNewUrl string = "/new"

// API - Users
var ApiUserUrl string = fmt.Sprintf("/user/{%s}", DashboardUserQueryId)

const ApiUserProfileUrl string = "/profile"

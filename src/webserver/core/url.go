package core

import "fmt"

// Landing Page
var GetStartedUrl string = "/getting-started"
var ContactUsUrl string = "/contact-us"
var HomePageUrl string = "/"
var LoginUrl string = "/login"
var LogoutUrl string = "/logout"
var LearnMoreUrl string = "/learn"
var SamlCallbackUrl string = LoadEnvConfig().Login.RedirectUrl

// Dashboard
var DashboardUrl string = "/dashboard"
var DashboardHomeUrl string = "/"

// Dashboard - Organization
var DashboardOrgOrgQueryId string = "orgId"
var DashboardOrgUrl string = fmt.Sprintf("/org/{%s}", DashboardOrgOrgQueryId)
var DashboardOrgHomeUrl string = "/"

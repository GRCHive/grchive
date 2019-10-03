package core

// Landing Page
var GetStartedUrl string = "/getting-started"
var ContactUsUrl string = "/contact-us"
var HomePageUrl string = "/"
var LoginUrl string = "/login"
var LearnMoreUrl string = "/learn"
var SamlCallbackUrl string = LoadEnvConfig().Login.RedirectUrl

// Dashboard
var DashboardUrl string = "/dashboard"
var DashboardHomeUrl string = "/"
var DashboardOrgHomeUrl string = "/org/{orgId}"

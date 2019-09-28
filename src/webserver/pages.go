package main

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
	"net/http"
)

func renderGettingStartedPage(w http.ResponseWriter, r *http.Request) {
	retrieveTemplate(GettingStartedPageTemplateKey).
		ExecuteTemplate(
			w,
			"base",
			core.LoadGlobalProps())
}

func renderContactUsPage(w http.ResponseWriter, r *http.Request) {
	retrieveTemplate(ContactUsPageTemplateKey).
		ExecuteTemplate(
			w,
			"base",
			core.LoadGlobalProps())
}

func renderHomePage(w http.ResponseWriter, r *http.Request) {
	retrieveTemplate(LandingPageTemplateKey).
		ExecuteTemplate(
			w,
			"base",
			core.LoadGlobalProps())
}

func renderLoginPage(w http.ResponseWriter, r *http.Request) {
	retrieveTemplate(LoginPageTemplateKey).
		ExecuteTemplate(
			w,
			"base",
			core.LoadGlobalProps())
}

func renderLearnMorePage(w http.ResponseWriter, r *http.Request) {
	retrieveTemplate(LearnMorePageTemplateKey).
		ExecuteTemplate(
			w,
			"base",
			core.LoadGlobalProps())
}

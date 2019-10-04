package render

import (
	"net/http"
)

func Render404(w http.ResponseWriter, r *http.Request) {
	RetrieveTemplate(Error404TemplateKey).
		ExecuteTemplate(
			w,
			"base",
			BuildTemplateParams(w, r, false))
}

func Render403(w http.ResponseWriter, r *http.Request) {
	RetrieveTemplate(Error403TemplateKey).
		ExecuteTemplate(
			w,
			"base",
			BuildTemplateParams(w, r, false))
}

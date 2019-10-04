package render

import (
	"net/http"
)

func Render404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	RetrieveTemplate(Error404TemplateKey).
		ExecuteTemplate(
			w,
			"base",
			BuildTemplateParams(w, r, false))
}

func Render403(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	RetrieveTemplate(Error403TemplateKey).
		ExecuteTemplate(
			w,
			"base",
			BuildTemplateParams(w, r, false))
}

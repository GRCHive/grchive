package render

import (
	"net/http"
)

func Render404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	RenderTemplate(w, Error404TemplateKey, "base", BuildPageTemplateParametersFull(r), emptyParams)
}

func Render403(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	RenderTemplate(w, Error403TemplateKey, "base", BuildPageTemplateParametersFull(r), emptyParams)
}

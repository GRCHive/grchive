package webcore

import "net/http"

// This response writer will change the original http.ResponseWriter so that
// it writes a status code of 301 -> 308 and 302 -> 307.
type RedirectResponseWriter struct {
	w http.ResponseWriter
}

func (this RedirectResponseWriter) Header() http.Header {
	return this.w.Header()
}

func (this RedirectResponseWriter) Write(b []byte) (int, error) {
	return this.w.Write(b)
}

func (this RedirectResponseWriter) WriteHeader(statusCode int) {
	if statusCode == http.StatusMovedPermanently {
		this.w.WriteHeader(http.StatusPermanentRedirect)
	} else if statusCode == http.StatusFound {
		this.w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		this.w.WriteHeader(statusCode)
	}
}

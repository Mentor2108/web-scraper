package rest

import (
	"backend-service/defn"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func ApplyMiddleware(router *httprouter.Router) http.Handler {
	return responseContentTypeJSON(router)
}

func responseContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(defn.HTTPHeaderContentType, defn.ContentTypeJSON)

		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

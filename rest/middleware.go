package rest

import (
	"backend-service/defn"
	"backend-service/util"
	"net/http"
	"runtime/debug"

	"github.com/julienschmidt/httprouter"
)

func ApplyMiddleware(router *httprouter.Router) http.Handler {
	return panicHandler(responseContentTypeJSON(router))
}

func responseContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(defn.HTTPHeaderContentType, defn.ContentTypeJSON)

		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

func panicHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			log := util.GetGlobalLogger(r.Context())
			if panicVal := recover(); panicVal != nil {
				log.Printf("Recovered in middleware:\n%+v\n%s\n", panicVal, string(debug.Stack()))
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "something went wrong in server's side"}`))
			}
		}()

		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

package defn

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request, params httprouter.Params)
}

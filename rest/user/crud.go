package user

import (
	"backend-service/util"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (h *UserRoutesHandler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	log := util.GetGlobalLogger(ctx)
	requestMap, cerr := util.ConvertRequestBodyToMap(ctx, r.Body)
	if cerr != nil {
		log.Println("Failed to parse request body", cerr)
		util.RespondWithError(ctx, w, http.StatusBadRequest, cerr)
		return
	}

	log.Println("Request map: ", requestMap)
	h.service.Create(ctx, requestMap)
}

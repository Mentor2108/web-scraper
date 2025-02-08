package rest

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func GetFileById(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	// ctx := r.Context()
	// log := util.GetGlobalLogger(ctx)

	// requestMap, cerr := util.ConvertRequestBodyToMap(ctx, r.Body)
	// if cerr != nil {
	// 	log.Println(cerr)
	// 	util.RespondWithError(ctx, w, http.StatusBadRequest, cerr)
	// 	return
	// }

	// fileId := param.ByName("id")
	// if strings.EqualFold(fileId, "") {
	// 	cerr := util.NewCustomError(ctx, "empty-file-id", errors.New("no file id provided"))
	// 	log.Println(cerr)
	// 	util.RespondWithError(ctx, w, http.StatusBadRequest, cerr)
	// 	return
	// }

	// service.GetFile(ctx)
}

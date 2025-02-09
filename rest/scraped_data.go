package rest

import (
	"backend-service/defn"
	"backend-service/service"
	"backend-service/util"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func GetScrapedDataForJobs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	log := util.GetGlobalLogger(ctx)

	queryParams := r.URL.Query()

	jobId := queryParams.Get("id")
	pagesize := queryParams.Get("pagesize")

	pagesizeInt, err := strconv.Atoi(pagesize)
	if err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeParseToIntFailed, defn.ErrParseToIntFailed, map[string]string{
			"field": "pagesize",
			"error": err.Error(),
		})
		log.Println(cerr)
		util.RespondWithError(ctx, w, http.StatusBadRequest, cerr)
		return
	}

	if resp, cerr := service.GetScrapeTasksForScrapeJob(ctx, jobId, pagesizeInt); cerr != nil {
		log.Println(cerr)
		util.RespondWithError(ctx, w, http.StatusBadRequest, cerr)
		return
	} else {
		util.SendResponseMapWithStatus(ctx, w, http.StatusOK, resp)
		return
	}
}

package rest

import (
	"backend-service/defn"
	"backend-service/service"
	"backend-service/util"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func ScrapeURL(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	log := util.GetGlobalLogger(ctx)

	var scrapeRequest defn.ScrapeRequest
	scrapeRequest.Config = &defn.ScrapeConfig{Depth: 5}
	if err := json.NewDecoder(r.Body).Decode(&scrapeRequest); err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeFailedToParseRequestBody, defn.ErrFailedToParseRequestBody, map[string]string{
			"error": err.Error(),
		})
		log.Println(cerr)
		util.RespondWithError(ctx, w, http.StatusBadRequest, cerr)
		return
	}

	// if val, ok := requestMap["config"]; !ok {
	if scrapeRequest.Config == nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeMissingRequiredField, defn.ErrMissingRequiredField, map[string]string{
			"field": "config",
		})
		log.Println(cerr)
		util.RespondWithError(ctx, w, http.StatusBadRequest, cerr)
		return
	}

	var urlScraper *service.UrlScraperService
	scraperService, cerr := urlScraper.Init(ctx, *scrapeRequest.Config, map[string]interface{}{"url": scrapeRequest.Url})
	if cerr != nil {
		log.Println(cerr)
		util.RespondWithError(ctx, w, http.StatusBadRequest, cerr)
		return
	}

	if resp, cerr := scraperService.Start(ctx); cerr != nil {
		log.Println(cerr)
		util.RespondWithError(ctx, w, http.StatusBadRequest, cerr)
		return
	} else {
		util.SendResponseMapWithStatus(ctx, w, http.StatusCreated, resp)
		return
	}
	// requestMap[""]

	// log.Printf("Request map: %+v", scrapeRequest)
	// spew.Dump(scrapeRequest)
	// h.service.Create(ctx, requestMap)
}

func ScrapePDF(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	cerr := util.NewCustomError(ctx, "not-implemented", errors.New("PDF scraper is not yet implemented"))
	util.GetGlobalLogger(ctx).Println(cerr)

	util.RespondWithError(ctx, w, http.StatusNotAcceptable, cerr)
}

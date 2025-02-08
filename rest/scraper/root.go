package scraper

import (
	"backend-service/defn"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ScraperRoutesHandler struct {
	service defn.Service
}

func NewScraperRoutesHandler(service defn.Service) *ScraperRoutesHandler {
	return &ScraperRoutesHandler{
		service: service,
	}
}

func AddRoutes(router *httprouter.Router, scraperHandler *ScraperRoutesHandler) {
	router.Handle(http.MethodGet, "/scraper/scrape/website/start", scraperHandler.Scrape)
}

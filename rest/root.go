package rest

import (
	"backend-service/defn"
	"backend-service/util"
	"context"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func AddRoutes(router *httprouter.Router) {
	log := util.GetGlobalLogger(context.Background())

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(time.Now().Format(time.RFC822Z), " 404 for - ", r.URL.String())
		w.WriteHeader(http.StatusNotFound)
	})

	router.Handle(http.MethodGet, "/status", ServerStatus)

	router.POST("/scraper/scrape/url/start", ScrapeURL)
	router.POST("/scraper/scrape/pdf/start", ScrapePDF)

	router.GET("/scraper/list/scrapeddata/url", GetScrapedDataForJobs)

	router.GET("/content/file/id/:id", GetFileById)
}

func ServerStatus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set(defn.HTTPHeaderContentType, defn.ContentTypePlainText)
	w.Write([]byte("Server is up and running"))
}

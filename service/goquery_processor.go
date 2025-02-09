package service

import (
	"backend-service/data"
	"backend-service/defn"
	"backend-service/util"
	"context"
)

type GoqueryProcessorService struct {
	config         defn.ScrapeConfig
	scrapeInfo     map[string]interface{}
	ScrapeJobRepo  *data.ScrapeJobRepo
	ScrapeTaskRepo *data.ScrapeTaskRepo
}

func (processor *GoqueryProcessorService) ProcessPhase(ctx context.Context, rawHtml string) (map[string]interface{}, *util.CustomError) {
	// scraper.config.ProcessPhase
	return nil, nil
}

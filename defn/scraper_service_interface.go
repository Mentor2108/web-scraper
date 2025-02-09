package defn

import (
	"backend-service/util"
	"context"
)

type ScraperService interface {
	Init(ctx context.Context, config ScrapeConfig, scrapeInfo map[string]interface{}) (ScraperService, *util.CustomError)
	Start(ctx context.Context) (map[string]interface{}, *util.CustomError)
	Pause(ctx context.Context) *util.CustomError
	Stop(ctx context.Context) *util.CustomError
	Status(ctx context.Context) (map[string]interface{}, *util.CustomError)
	SyncStart(ctx context.Context) (map[string]interface{}, *util.CustomError)
}

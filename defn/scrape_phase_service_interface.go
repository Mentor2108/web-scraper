package defn

import (
	"backend-service/util"
	"context"
)

type ScrapePhaseService interface {
	Init(ctx context.Context, config ScrapeConfig, scrapeInfo map[string]interface{}) (ScrapePhaseService, *util.CustomError)
	Start(ctx context.Context) (string, map[string]interface{}, *util.CustomError)
	Pause(ctx context.Context) *util.CustomError
	Stop(ctx context.Context) *util.CustomError
	Status(ctx context.Context) (map[string]interface{}, *util.CustomError)
}

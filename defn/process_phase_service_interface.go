package defn

import (
	"backend-service/util"
	"context"
)

type ProcessPhaseService interface {
	Init(ctx context.Context, config ScrapeConfig, scrapeInfo map[string]interface{}) (ProcessPhaseService, *util.CustomError)
	Process(ctx context.Context, rawHTML string) (string, map[string]interface{}, *util.CustomError)
}

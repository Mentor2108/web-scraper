package service

import (
	"backend-service/data"
	"backend-service/util"
	"context"
	"strings"
)

func GetScrapeTasksForScrapeJob(ctx context.Context, jobId string, pagesize int) (map[string]interface{}, *util.CustomError) {
	scrapeJobRepo := data.NewScrapeJobRepo()
	if strings.EqualFold(jobId, "") {
		return scrapeJobRepo.GetJobWithTasks(ctx, pagesize)
	} else {
		return scrapeJobRepo.GetJobWithTasksByID(ctx, jobId)
	}
}

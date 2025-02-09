package service

import (
	"backend-service/data"
	"backend-service/defn"
	"backend-service/util"
	"context"
)

type UrlScraperService struct {
	scraper        defn.SpecialisedScraperService
	config         defn.ScrapeConfig
	scrapeInfo     map[string]interface{}
	ScrapeJobRepo  *data.ScrapeJobRepo
	ScrapeTaskRepo *data.ScrapeTaskRepo
}

// func NewScraperService(repo defn.Repository) *ScraperService {
// 	return &ScraperService{scraper: repo}
// }

func (scraper *UrlScraperService) Init(ctx context.Context, config defn.ScrapeConfig, scrapeInfo map[string]interface{}) (defn.ScraperService, *util.CustomError) {
	log := util.GetGlobalLogger(ctx)
	var cerr *util.CustomError
	var specialisedScraper defn.SpecialisedScraperService

	switch config.ScrapePhase.Library {
	case defn.ScrapePhaseLibraryChromedp:
		var ChromedpScraper *ChromedpScraperService
		if specialisedScraper, cerr = ChromedpScraper.Init(ctx, config, scrapeInfo); cerr != nil {
			return nil, cerr
		}
	default:
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeScrapePhaseLibraryNotSupported, defn.ErrScrapePhaseLibraryNotSupported, map[string]string{
			"library": config.ScrapePhase.Library,
		})
		log.Println(cerr)
		return nil, cerr
	}

	switch config.ProcessPhase.Library {
	case defn.ProcessPhaseLibraryGoquery:
		var GoqueryProcessor
	}

	return &UrlScraperService{
		scraper:        specialisedScraper,
		config:         config,
		scrapeInfo:     scrapeInfo,
		ScrapeJobRepo:  data.NewScrapeJobRepo(),
		ScrapeTaskRepo: data.NewScrapeTaskRepo(),
	}, nil
}

func (scraper *UrlScraperService) Start(ctx context.Context) (map[string]interface{}, *util.CustomError) {
	//call a go subroutine
	go func() {
		//Currently not storing any extra info in the context so creating a new context is fine
		//Otherwise would have to copy the needed keys from request context to the new one
		scraper.SyncStart(context.Background())
	}()

	//exit with the generated id
	return map[string]interface{}{
		"response": map[string]interface{}{
			"status": "successfully started scraping",
			"job_id": scraper.scrapeInfo["job-id"].(string),
		},
	}, nil
}

func (scraper *UrlScraperService) Pause(ctx context.Context) *util.CustomError {
	return nil
}

func (scraper *UrlScraperService) Stop(ctx context.Context) *util.CustomError {
	return nil
}

func (scraper *UrlScraperService) Status(ctx context.Context) (map[string]interface{}, *util.CustomError) {
	return nil, nil
}

func (scraper *UrlScraperService) SyncStart(ctx context.Context) (map[string]interface{}, *util.CustomError) {
	//create scrape-job
	// log := util.GetGlobalLogger(ctx)

	// jobId, cerr := scraper.ScrapeJobRepo.Create(ctx, defn.ScrapeJob{
	// 	URL:      scraper.scrapeInfo["url"].(string),
	// 	Depth:    scraper.config.Depth,
	// 	Maxlimit: scraper.config.MaxLimit,
	// })
	// if cerr != nil {
	// 	log.Println(cerr)
	// 	return nil, cerr
	// }

	// taskId, cerr := scraper.ScrapeTaskRepo.Create(ctx, defn.ScrapeTask{
	// 	URL:      scraper.scrapeInfo["url"].(string),
	// 	JobId:    jobId,
	// 	Depth:    scraper.config.Depth,
	// 	Maxlimit: scraper.config.MaxLimit,
	// 	Level:    1,
	// })
	// if cerr != nil {
	// 	log.Println(cerr)
	// 	return nil, cerr
	// }

	//create scrape-link-obj

	// scraper.scrapeInfo["job-id"] = jobId
	// scraper.scrapeInfo["task-id"] = taskId

	return scraper.scraper.Start(context.Background())
}

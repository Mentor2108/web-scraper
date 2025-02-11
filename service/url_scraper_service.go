package service

import (
	"backend-service/data"
	"backend-service/defn"
	"backend-service/util"
	"context"
	"errors"
	"runtime/debug"
	"strings"
)

type UrlScraperService struct {
	scrapePhase    defn.ScrapePhaseService
	processPhase   defn.ProcessPhaseService
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
	var scrapePhaseScraper defn.ScrapePhaseService
	var processPhaseProcessor defn.ProcessPhaseService

	if strings.EqualFold(config.Root, "") {
		cerr := util.NewCustomError(ctx, "empty-root-selector", errors.New("provided root selector is empty"))
		log.Println(cerr)
		return nil, cerr
	}

	switch config.ScrapePhase.Library {
	case defn.ScrapePhaseLibraryChromedp:
		var ChromedpScraper *ChromedpScraperService
		if scrapePhaseScraper, cerr = ChromedpScraper.Init(ctx, config, scrapeInfo); cerr != nil {
			return nil, cerr
		}
	default:
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeScrapePhaseLibraryNotSupported, defn.ErrScrapePhaseLibraryNotSupported, map[string]string{
			"library": config.ScrapePhase.Library,
		})
		log.Println(cerr)
		return nil, cerr
	}

	if config.ProcessPhase != nil {
		switch config.ProcessPhase.Library {
		case defn.ProcessPhaseLibraryGoquery:
			var GoqueryProcessor *GoqueryProcessURL
			if processPhaseProcessor, cerr = GoqueryProcessor.Init(ctx, config, scrapeInfo); cerr != nil {
				return nil, cerr
			}
		default:
			cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeProcessPhaseLibraryNotSupported, defn.ErrProcessPhaseLibraryNotSupported, map[string]string{
				"library": config.ProcessPhase.Library,
			})
			log.Println(cerr)
			return nil, cerr
		}
	}

	return &UrlScraperService{
		scrapePhase:    scrapePhaseScraper,
		processPhase:   processPhaseProcessor,
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
		defer func() {
			log := util.GetGlobalLogger(context.Background())
			if panicVal := recover(); panicVal != nil {
				log.Printf("Recovered in middleware:\n%+v\n%s\n", panicVal, string(debug.Stack()))
			}
		}()
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
	log := util.GetGlobalLogger(ctx)

	rawHtml, resp, cerr := scraper.scrapePhase.Start(context.Background())
	if cerr != nil {
		// if !scraper.config.ContinueOnError || strings.EqualFold(rawHtml, "") {
		//write error to database for scrape job
		errorResponse, databaseErr := scraper.ScrapeJobRepo.Update(ctx, scraper.scrapeInfo["job-id"].(string), map[string]interface{}{
			"response": map[string]interface{}{
				"status":         "scraping failed at task: " + scraper.scrapeInfo["task-id"].(string),
				"error":          cerr.GetErrorMap(ctx),
				"uploaded_files": scraper.scrapeInfo["all_uploaded_files"],
			},
		})
		if databaseErr != nil {
			log.Println(databaseErr)
			return resp, databaseErr
		}
		log.Println("error response:", errorResponse)
		return resp, cerr
		// }
	}

	var returnedConfig map[string]interface{} = resp
	if scraper.config.ProcessPhase != nil {
		_, returnedConfig, cerr = scraper.processPhase.Process(ctx, rawHtml)
		// scraper.scrapeInfo["uploaded_files"] = append((scraper.scrapeInfo["uploaded_files"].([]map[string]interface{})), (returnedConfig["uploaded_files"].([]map[string]interface{}))...)

		if cerr != nil {
			errorResponse, databaseErr := scraper.ScrapeJobRepo.Update(ctx, scraper.scrapeInfo["job-id"].(string), map[string]interface{}{
				"response": map[string]interface{}{
					"status":         "scraping failed at task: " + scraper.scrapeInfo["task-id"].(string),
					"error":          cerr.GetErrorMap(ctx),
					"uploaded_files": scraper.scrapeInfo["all_uploaded_files"],
				},
			})
			if databaseErr != nil {
				log.Println(databaseErr)
				return returnedConfig, databaseErr
			}
			log.Println("error response:", errorResponse)
			return returnedConfig, cerr
		}
	}

	_, databaseErr := scraper.ScrapeJobRepo.Update(ctx, scraper.scrapeInfo["job-id"].(string), map[string]interface{}{
		"response": map[string]interface{}{
			"status":         "successfully scraped provided url",
			"uploaded_files": scraper.scrapeInfo["all_uploaded_files"],
		},
	})
	if databaseErr != nil {
		return returnedConfig, databaseErr
	}

	return returnedConfig, nil
}

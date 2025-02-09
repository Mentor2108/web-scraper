package service

import (
	"backend-service/data"
	"backend-service/defn"
	"backend-service/util"
	"context"
	"log"

	"github.com/chromedp/chromedp"
)

type ChromedpScraperService struct {
	config             defn.ScrapeConfig
	scrapeInfo         map[string]interface{}
	chromedpContext    context.Context
	chromedpCancelFunc context.CancelFunc
	ScrapeJobRepo      *data.ScrapeJobRepo
	ScrapeTaskRepo     *data.ScrapeTaskRepo
}

// func NewChromedpScraperService(repo defn.Repository) *ScraperService {
// 	return &ScraperService{repo: repo}
// }

func (scraper *ChromedpScraperService) Init(ctx context.Context, config defn.ScrapeConfig, scrapeInfo map[string]interface{}) (defn.SpecialisedScraperService, *util.CustomError) {
	// log := util.GetGlobalLogger(ctx)
	newCtx, cancelFunc := chromedp.NewContext(context.Background())

	specalisedScraper := &ChromedpScraperService{
		config:             config,
		scrapeInfo:         scrapeInfo,
		chromedpContext:    newCtx,
		chromedpCancelFunc: cancelFunc,
		ScrapeJobRepo:      data.NewScrapeJobRepo(),
		ScrapeTaskRepo:     data.NewScrapeTaskRepo(),
	}

	jobId, cerr := specalisedScraper.ScrapeJobRepo.Create(ctx, defn.ScrapeJob{
		Depth:    specalisedScraper.config.Depth,
		Maxlimit: specalisedScraper.config.MaxLimit,
	})
	if cerr != nil {
		log.Println(cerr)
		return nil, cerr
	}
	specalisedScraper.scrapeInfo["job-id"] = jobId

	log.Println("jobid:", jobId)
	taskId, cerr := specalisedScraper.ScrapeTaskRepo.Create(ctx, defn.ScrapeTask{
		URL:      specalisedScraper.scrapeInfo["url"].(string),
		JobId:    jobId,
		Depth:    specalisedScraper.config.Depth,
		Maxlimit: specalisedScraper.config.MaxLimit,
		Level:    1,
	})
	if cerr != nil {
		log.Println(cerr)
		return nil, cerr
	}
	specalisedScraper.scrapeInfo["task-id"] = taskId

	return specalisedScraper, nil
}

func (scraper *ChromedpScraperService) Start(ctx context.Context) (map[string]interface{}, *util.CustomError) {
	log := util.GetGlobalLogger(ctx)
	rawHTML, plainText, resp, cerr := scraper.start(ctx)
	if cerr != nil {
		//write error to database
		errorResponse, databaseErr := scraper.ScrapeTaskRepo.Update(ctx, scraper.scrapeInfo["task-id"].(string), map[string]interface{}{
			"response": map[string]interface{}{
				"status": "scraping failed",
				"error":  cerr.GetErrorMap(ctx),
			},
		})
		if databaseErr != nil {
			log.Println(databaseErr)
			return resp, databaseErr
		}
		log.Println("error response:", errorResponse)
	}
	return resp, cerr
}

func (scraper *ChromedpScraperService) Pause(ctx context.Context) *util.CustomError {
	return nil
}

func (scraper *ChromedpScraperService) Stop(ctx context.Context) *util.CustomError {
	return nil
}

func (scraper *ChromedpScraperService) Status(ctx context.Context) (map[string]interface{}, *util.CustomError) {
	return nil, nil
}

func (scraper *ChromedpScraperService) start(ctx context.Context) (string, string, map[string]interface{}, *util.CustomError) {
	log := util.GetGlobalLogger(ctx)

	//saved it in the map, can assert its there
	url := scraper.scrapeInfo["url"].(string)

	var rawHTML string
	var plainText string

	defer scraper.chromedpCancelFunc()
	err := chromedp.Run(scraper.chromedpContext, chromedp.ActionFunc(func(chromedpCtx context.Context) error {
		if err := chromedp.Navigate(url).Do(chromedpCtx); err != nil {
			// log.Println(err)
			return err
		}

		if err := chromedp.Text(scraper.config.Root, &plainText).Do(chromedpCtx); err != nil {
			return err
		}

		if err := chromedp.OuterHTML(scraper.config.Root, &rawHTML).Do(chromedpCtx); err != nil {
			return err
		}
		return nil
	}))
	if err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeChromedpError, defn.ErrChromepError, map[string]string{
			"error": err.Error(),
		})
		log.Println(cerr)
		return "", "", nil, cerr
	}

	// change job-id to task-id later
	savedFiles, cerr := ParseFolderStructureAndSaveFile(ctx, "static", &defn.FileFolderStructure{
		Name: "scraped_data",
		Folders: []*defn.FileFolderStructure{
			{
				Name: scraper.scrapeInfo["task-id"].(string),
				Files: []*defn.FileStructure{
					{
						FileName:    "raw_html",
						FileType:    ".html",
						FileContent: []byte(rawHTML),
					},
					{
						FileName:    "plain_text",
						FileType:    ".txt",
						FileContent: []byte(plainText),
					},
				},
			},
		},
	})
	uploadedFilesResponse := map[string]interface{}{
		"uploaded_files": savedFiles,
	}
	if cerr != nil {
		log.Println(cerr)
		return rawHTML, plainText, uploadedFilesResponse, cerr
	}

	updateResponse, cerr := scraper.ScrapeTaskRepo.Update(ctx, scraper.scrapeInfo["task-id"].(string), map[string]interface{}{
		"response": map[string]interface{}{
			"status":         "successfully scraped provided url",
			"uploaded_files": savedFiles,
		},
	})
	if cerr != nil {
		return rawHTML, plainText, uploadedFilesResponse, cerr
	}

	_, cerr = scraper.ScrapeJobRepo.Update(ctx, scraper.scrapeInfo["job-id"].(string), map[string]interface{}{
		"response": map[string]interface{}{
			"status":         "successfully scraped provided url",
			"uploaded_files": savedFiles,
		},
	})
	if cerr != nil {
		return rawHTML, plainText, uploadedFilesResponse, cerr
	}

	log.Println("sucess scraping response:", updateResponse)
	return rawHTML, plainText, uploadedFilesResponse, nil
}

package service

import (
	"backend-service/data"
	"backend-service/defn"
	"backend-service/util"
	"context"

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

func (scraper *ChromedpScraperService) Init(ctx context.Context, config defn.ScrapeConfig, scrapeInfo map[string]interface{}) (defn.ScrapePhaseService, *util.CustomError) {
	log := util.GetGlobalLogger(ctx)
	newCtx, _ := chromedp.NewContext(context.Background())
	newCtxWithTimeout, cancelFuncWithTimeout := context.WithTimeout(newCtx, defn.ChromedpTimeout)

	specalisedScraper := &ChromedpScraperService{
		config:             config,
		scrapeInfo:         scrapeInfo,
		chromedpContext:    newCtxWithTimeout,
		chromedpCancelFunc: cancelFuncWithTimeout,
		ScrapeJobRepo:      data.NewScrapeJobRepo(),
		ScrapeTaskRepo:     data.NewScrapeTaskRepo(),
	}

	var jobId string
	var cerr *util.CustomError
	if _, ok := specalisedScraper.scrapeInfo["job-id"]; !ok { //First run
		jobId, cerr = specalisedScraper.ScrapeJobRepo.Create(ctx, defn.ScrapeJob{
			Depth:    specalisedScraper.config.Depth,
			Maxlimit: specalisedScraper.config.MaxLimit,
		})
		if cerr != nil {
			log.Println(cerr)
			return nil, cerr
		}
		specalisedScraper.scrapeInfo["job-id"] = jobId
		specalisedScraper.scrapeInfo["all_uploaded_files"] = map[string]interface{}{}
		scrapeInfo["visitedurls"] = []string{scrapeInfo["url"].(string)}
	} else {
		jobId = scrapeInfo["job-id"].(string)
	}

	taskId, cerr := specalisedScraper.ScrapeTaskRepo.Create(ctx, defn.ScrapeTask{
		URL:      specalisedScraper.scrapeInfo["url"].(string),
		JobId:    jobId,
		Depth:    specalisedScraper.config.Depth,
		Maxlimit: specalisedScraper.config.MaxLimit,
		Level:    scrapeInfo["level"].(int),
	})
	if cerr != nil {
		log.Println(cerr)
		return nil, cerr
	}
	specalisedScraper.scrapeInfo["task-id"] = taskId

	return specalisedScraper, nil
}

func (scraper *ChromedpScraperService) Start(ctx context.Context) (string, map[string]interface{}, *util.CustomError) {
	log := util.GetGlobalLogger(ctx)
	rawHTML, resp, cerr := scraper.start(ctx)
	if cerr != nil {
		//write error to database
		errorResponse, databaseErr := scraper.ScrapeTaskRepo.Update(ctx, scraper.scrapeInfo["task-id"].(string), map[string]interface{}{
			"response": map[string]interface{}{
				"status":         "scraping failed at scrape phase",
				"error":          cerr.GetErrorMap(ctx),
				"uploaded_files": scraper.scrapeInfo["uploaded_files"],
			},
		})
		if databaseErr != nil {
			log.Println(databaseErr)
			return rawHTML, resp, databaseErr
		}
		log.Println("error response:", errorResponse)
	}
	return rawHTML, resp, cerr
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

func (scraper *ChromedpScraperService) start(ctx context.Context) (string, map[string]interface{}, *util.CustomError) {
	log := util.GetGlobalLogger(ctx)

	//saved it in the map, can assert its there
	url := scraper.scrapeInfo["url"].(string)

	var rawHTML string
	// var plainText string

	defer scraper.chromedpCancelFunc()
	err := chromedp.Run(scraper.chromedpContext, chromedp.ActionFunc(func(chromedpCtx context.Context) error {
		if err := chromedp.Navigate(url).Do(chromedpCtx); err != nil {
			// log.Println(err)
			return err
		}

		// if err := chromedp.Text(scraper.config.Root, &plainText).Do(chromedpCtx); err != nil {
		// 	return err
		// }

		if err := chromedp.OuterHTML(scraper.config.Root, &rawHTML).Do(chromedpCtx); err != nil {
			return err
		}
		return nil
	}))
	if err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeChromedpError, defn.ErrChromedpError, map[string]string{
			"error": err.Error(),
		})
		log.Println(cerr)
		return "", nil, cerr
	}

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
					// {
					// 	FileName:    "plain_text",
					// 	FileType:    ".md",
					// 	FileContent: []byte(plainText),
					// },
				},
			},
		},
	})
	uploadedFilesResponse := map[string]interface{}{
		"uploaded_files": savedFiles,
	}
	scraper.scrapeInfo["uploaded_files"] = savedFiles

	scraper.scrapeInfo["all_uploaded_files"].(map[string]interface{})[scraper.scrapeInfo["task-id"].(string)] = scraper.scrapeInfo["uploaded_files"]

	if cerr != nil {
		log.Println(cerr)
		return rawHTML, uploadedFilesResponse, cerr
	}

	_, cerr = scraper.ScrapeTaskRepo.Update(ctx, scraper.scrapeInfo["task-id"].(string), map[string]interface{}{
		"response": map[string]interface{}{
			"status":         "successfully scraped provided url",
			"uploaded_files": savedFiles,
		},
	})
	if cerr != nil {
		return rawHTML, uploadedFilesResponse, cerr
	}

	// log.Println("success scraping response:", updateResponse)
	return rawHTML, uploadedFilesResponse, nil
}

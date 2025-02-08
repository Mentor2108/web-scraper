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
	// ScrapeTaskRepo data.ScrapeTaskRepo
}

// func NewChromedpScraperService(repo defn.Repository) *ScraperService {
// 	return &ScraperService{repo: repo}
// }

func (scraper *ChromedpScraperService) Init(ctx context.Context, config defn.ScrapeConfig, scrapeInfo map[string]interface{}) (defn.SpecialisedScraperService, *util.CustomError) {
	// log := util.GetGlobalLogger(ctx)
	newCtx, cancelFunc := chromedp.NewContext(context.Background())

	return &ChromedpScraperService{
		config:             config,
		scrapeInfo:         scrapeInfo,
		chromedpContext:    newCtx,
		chromedpCancelFunc: cancelFunc,
		ScrapeJobRepo:      data.NewScrapeJobRepo(),
	}, nil
}

func (scraper *ChromedpScraperService) Start(ctx context.Context) *util.CustomError {
	log := util.GetGlobalLogger(ctx)
	cerr := scraper.start(ctx)
	if cerr != nil {
		//write error to database
		errorResponse, databaseErr := scraper.ScrapeJobRepo.Update(ctx, scraper.scrapeInfo["job-id"].(string), map[string]interface{}{
			"response": map[string]interface{}{
				"status": "scraping failed",
				"error":  cerr.GetErrorMap(ctx),
			},
		})
		if databaseErr != nil {
			log.Println(databaseErr)
			return databaseErr
		}
		log.Println("error response:", errorResponse)
	}
	return cerr
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

func (scraper *ChromedpScraperService) start(ctx context.Context) *util.CustomError {
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
		return cerr
	}

	// change job-id to task-id later
	savedFiles, cerr := ParseFolderStructureAndSaveFile(ctx, "static", &defn.FileFolderStructure{
		Name: "scraped_data",
		Folders: []*defn.FileFolderStructure{
			{
				Name: scraper.scrapeInfo["job-id"].(string),
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
	if cerr != nil {
		log.Println(cerr)
		return cerr
	}

	updateResponse, cerr := scraper.ScrapeJobRepo.Update(ctx, scraper.scrapeInfo["job-id"].(string), map[string]interface{}{
		"response": map[string]interface{}{
			"status":         "successfully scraped provided url",
			"uploaded_files": savedFiles,
		},
	})
	if cerr != nil {
		return cerr
	}
	log.Println("sucess scraping response:", updateResponse)
	//save raw html somewhere
	// log.Println("plain text:", plainText)
	return nil
}

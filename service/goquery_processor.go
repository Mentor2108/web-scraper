package service

import (
	"backend-service/data"
	"backend-service/defn"
	"backend-service/util"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type GoqueryProcessURL struct {
	config         defn.ScrapeConfig
	scrapeInfo     map[string]interface{}
	fileIDs        []map[string]interface{}
	visitedLinks   []string
	ScrapeJobRepo  *data.ScrapeJobRepo
	ScrapeTaskRepo *data.ScrapeTaskRepo
}

func (processor *GoqueryProcessURL) Init(ctx context.Context, config defn.ScrapeConfig, scrapeInfo map[string]interface{}) (defn.ProcessPhaseService, *util.CustomError) {
	// processor.visitedLinks =
	// processor.fileIDs = []map[string]interface{}{}

	return &GoqueryProcessURL{
		config:         config,
		scrapeInfo:     scrapeInfo,
		ScrapeJobRepo:  data.NewScrapeJobRepo(),
		ScrapeTaskRepo: data.NewScrapeTaskRepo(),
		fileIDs:        []map[string]interface{}{},
		visitedLinks:   scrapeInfo["visitedurls"].([]string),
	}, nil
}
func (processor *GoqueryProcessURL) Process(ctx context.Context, rawHTML string) (string, map[string]interface{}, *util.CustomError) {
	log := util.GetGlobalLogger(ctx)
	processedStr, returnedConfig, cerr := processor.process(ctx, rawHTML)
	if cerr != nil {
		//write error to database
		errorResponse, databaseErr := processor.ScrapeTaskRepo.Update(ctx, processor.scrapeInfo["task-id"].(string), map[string]interface{}{
			"response": map[string]interface{}{
				"status":         "scraping failed at process phase",
				"error":          cerr.GetErrorMap(ctx),
				"uploaded_files": processor.scrapeInfo["uploaded_files"],
			},
		})
		if databaseErr != nil {
			log.Println(databaseErr)
			return processedStr, returnedConfig, databaseErr
		}
		log.Println("error response:", errorResponse)
	}
	return processedStr, returnedConfig, cerr
}

func (processor *GoqueryProcessURL) process(ctx context.Context, rawHTML string) (string, map[string]interface{}, *util.CustomError) {
	log := util.GetGlobalLogger(ctx)
	var returnedConfig map[string]interface{}

	processedStr, linksFound, imagesFound, _, cerr := processor.parseHTMLUsingDefinition(ctx, rawHTML)
	if cerr != nil {
		return processedStr, nil, cerr
	}
	ext := ".md"

	var imagesFileStruct []*defn.FileStructure
	for _, imageMap := range imagesFound {
		imagesFileStruct = append(imagesFileStruct, &defn.FileStructure{
			FileName:    imageMap["name"].(string),
			FileType:    imageMap["extension"].(string),
			FileContent: imageMap["content"].([]byte),
		})
	}

	filesMap, cerr := ParseFolderStructureAndSaveFile(ctx, "static", &defn.FileFolderStructure{
		Name: "scraped_data",
		Folders: []*defn.FileFolderStructure{
			{
				Name: processor.scrapeInfo["task-id"].(string),
				Files: []*defn.FileStructure{
					{
						FileName:    "plain_text",
						FileType:    ext,
						FileContent: []byte(processedStr),
					},
				},
				Folders: []*defn.FileFolderStructure{
					{
						Name:  "images",
						Files: imagesFileStruct,
					},
				},
			},
		},
	})

	if cerr != nil {
		return processedStr, nil, cerr
	}

	processor.fileIDs = append(processor.fileIDs, filesMap...)

	// if cerr != nil && !processor.config.ContinueOnError {
	// 	return processedStr, nil, cerr
	// } else if cerr == nil {
	// 	processor.fileIDs = append(processor.fileIDs, filesMap...)
	// }

	// log.Println("my depth, limit and visitedlinks are", processor.config.Depth, processor.config.MaxLimit, processor.visitedLinks)
	returnedConfig, cerr = processor.nestedScrape(ctx, linksFound)
	if cerr != nil {
		return processedStr, nil, cerr
	}

	processor.scrapeInfo["uploaded_files"] = append((processor.scrapeInfo["uploaded_files"].([]map[string]interface{})), processor.fileIDs...)
	processor.scrapeInfo["all_uploaded_files"].(map[string]interface{})[processor.scrapeInfo["task-id"].(string)] = processor.scrapeInfo["uploaded_files"]

	tempResponseMap := map[string]interface{}{
		"response": map[string]interface{}{
			"status":         "successfully processed provided url",
			"uploaded_files": processor.scrapeInfo["uploaded_files"],
		},
	}

	if returnedConfig == nil {
		returnedConfig = map[string]interface{}{
			"response": map[string]interface{}{
				"status":         "successfully processed provided url",
				"uploaded_files": processor.scrapeInfo["uploaded_files"],
			},
			"visitedurls": processor.visitedLinks,
		}
	} else {
		returnedConfig["response"] = tempResponseMap["response"]
	}

	_, cerr = processor.ScrapeTaskRepo.Update(ctx, processor.scrapeInfo["task-id"].(string), tempResponseMap)
	if cerr != nil {
		log.Println(cerr)
		return processedStr, returnedConfig, cerr
	}

	return processedStr, returnedConfig, nil
}

func (processor *GoqueryProcessURL) parseHTMLUsingDefinition(ctx context.Context, htmlStr string) (string, []string, []map[string]interface{}, *goquery.Document, *util.CustomError) {
	log := util.GetGlobalLogger(ctx)
	// config := processor.config
	// scrapingDefinitionName := config.Get("name").Str()
	processedHTML := ""

	var relevantLinks []string
	var cerr *util.CustomError

	htmlDoc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		cerr = util.NewCustomErrorWithKeys(ctx, defn.ErrCodeGoqueryError, defn.ErrGoqueryError, map[string]string{
			"error": err.Error(),
		})
		log.Println(cerr)
		return "", nil, nil, nil, cerr
	}

	// rootSelector := doc.Find(procesor.config.Root).First()
	// if rootSelector.Length() == 0 {
	// 	cerr := util.NewCustomError(ctx, defn.ErrCodeEmptyRootSelector, defn.ErrEmptyRootSelector)
	// 	log.Println(cerr)
	// 	return nil, cerr
	// }

	//globally excluding all elements
	excludeElements := processor.config.ExcludeElements
	for _, excludeElement := range excludeElements {
		htmlDoc.Find(excludeElement).Remove()
	}

	rootSelection := htmlDoc.Find(processor.config.Root).First() //applying first here to be sure there will always be one selection at this level
	if rootSelection.Length() < 1 {                              //if there's no match of 'root' selection, defaults to selecting body
		rootSelection = htmlDoc.Find("body").First()
	}

	var imagesData []map[string]interface{}
	scrapeImages := processor.config.ScrapeImages
	if scrapeImages {
		imagesData = processor.scrapeImage(ctx, rootSelection.Find("img"))
	}

	contentsToScrape := processor.config.ScrapeDataContent
	for idx, scrapingQuery := range contentsToScrape {
		prefix := "\n"
		suffix := "\n"

		scrapingQueryName := scrapingQuery.Name
		if strings.EqualFold(scrapingQueryName, "") {
			scrapingQueryName = fmt.Sprintf("scraping-query-%d", (idx + 1))
		}
		scrapingQueryType := scrapingQuery.Type

		selectedEle := rootSelection.Find(scrapingQuery.Selector)

		switch strings.ToLower(scrapingQueryType) {
		case defn.ScrapeQueryTypeText:
			if scrapingQuery.TextType != nil {
				prefix = scrapingQuery.TextType.Prefix
				suffix = scrapingQuery.TextType.Suffix
			}

			relevantLinks = append(relevantLinks, util.GetLinksFromSelection(ctx, selectedEle)...)
			processedHTML += prefix + util.TextWithoutSpaces(ctx, selectedEle) + suffix

		case defn.ScrapeQueryTypeTable:
			scrapedData, relevantTableLinks := processor.parseTableQueryType(ctx, scrapingQuery, selectedEle)
			relevantLinks = append(relevantLinks, relevantTableLinks...)
			processedHTML += scrapedData

		case defn.ScrapeQueryTypeSection:
			scrapedData, relevantSectionLinks := processor.parseSectionQueryType(ctx, scrapingQuery, selectedEle)
			relevantLinks = append(relevantLinks, relevantSectionLinks...)
			processedHTML += scrapedData

		default:
			log.Println("invalid type given for scrapingquery", scrapingQueryName)
			continue
		}

	}

	return processedHTML, relevantLinks, imagesData, htmlDoc, nil
}

func (processor *GoqueryProcessURL) parseTableQueryType(ctx context.Context, scrapingQuery defn.ScrapeDataContentDefn, rootSelection *goquery.Selection) (string, []string) {
	var processedStrBuilder strings.Builder
	var relevantLinks []string

	if scrapingQuery.TableType == nil {
		return "", nil
	}

	prefix := scrapingQuery.TableType.Prefix
	suffix := scrapingQuery.TableType.Suffix
	if strings.EqualFold(suffix, "") {
		suffix = "\n"
	}
	title := scrapingQuery.TableType.Title

	var columnLen int
	if scrapingQuery.TableType.ColumnsMap != nil {
		key := scrapingQuery.TableType.ColumnsMap.Key
		value := scrapingQuery.TableType.ColumnsMap.Value
		columnLen = 2
		rootSelection.Each(func(i int, EachSelectedTable *goquery.Selection) {
			relevantLinks = append(relevantLinks, util.GetLinksFromSelection(ctx, EachSelectedTable)...)
			processedStrBuilder.WriteString(prefix)

			if !strings.EqualFold(title, "") {
				text := EachSelectedTable.Find(title).Text()
				if !strings.EqualFold(text, "") {
					processedStrBuilder.WriteString("## " + text + "\n")
				}
			}
			processedStrBuilder.WriteString("| | |\n")
			processedStrBuilder.WriteString("| - | - |\n")

			EachSelectedTable.Find("tbody").Children().Each(func(i int, tableSelector *goquery.Selection) {
				if tableSelector.Is("tr") {
					noOfColumns := 0

					tableSelector.Children().Each(func(i int, rowSelector *goquery.Selection) {
						if rowSelector.Is(key) { //if its key
							processedStrBuilder.WriteString("|" + util.TextWithoutSpaces(ctx, rowSelector))
							noOfColumns++
						} else if rowSelector.Is(value) { //if its value
							processedStrBuilder.WriteString("|" + util.TextWithoutSpaces(ctx, rowSelector))
							noOfColumns++
						}
					})

					for idx := 0; idx < columnLen-noOfColumns; idx++ {
						processedStrBuilder.WriteString(" | ")
					}
					processedStrBuilder.WriteString("|\n")
				}
			})
			processedStrBuilder.WriteString(suffix)
		})
	} else if scrapingQuery.TableType.ColumnsNamesList != nil {
		// columnKey[]
		rootSelection.Each(func(i int, EachSelectedTable *goquery.Selection) {
			relevantLinks = append(relevantLinks, util.GetLinksFromSelection(ctx, EachSelectedTable)...)
			processedStrBuilder.WriteString(prefix)

			if !strings.EqualFold(title, "") {
				text := EachSelectedTable.Find(title).Text()
				if !strings.EqualFold(text, "") {
					processedStrBuilder.WriteString("## " + text + "\n")
				}
			}

			for columnIdx, columnName := range scrapingQuery.TableType.ColumnsNamesList {
				if columnIdx == 0 {
					processedStrBuilder.WriteString("|")
				}
				columnLen++

				if !strings.EqualFold(columnName, "") {
					processedStrBuilder.WriteString(columnName + "|")
				} else {
					processedStrBuilder.WriteString("column-" + fmt.Sprint(columnIdx) + "|")
				}
			}
			for columnIdx := range scrapingQuery.TableType.ColumnsNamesList {
				if columnIdx == 0 {
					processedStrBuilder.WriteString("\n|")
				}

				processedStrBuilder.WriteString(" - |")
			}
			processedStrBuilder.WriteString("\n")

			EachSelectedTable.Find("tbody").Children().Each(func(i int, tableSelector *goquery.Selection) {
				if tableSelector.Is("tr") {
					noOfColumns := 0
					// not handling if there are more columns (th/tb) in the row than no-of columns defined
					tableSelector.Children().Each(func(i int, rowSelector *goquery.Selection) {
						if noOfColumns == columnLen {
							return
						}

						processedStrBuilder.WriteString("|" + util.TextWithoutSpaces(ctx, rowSelector))
						noOfColumns++
					})

					for idx := 0; idx < columnLen-noOfColumns; idx++ {
						processedStrBuilder.WriteString(" | ")
					}
					processedStrBuilder.WriteString("|\n")
				}
			})
			processedStrBuilder.WriteString(suffix)
		})

	}
	return processedStrBuilder.String(), relevantLinks
}

func (processor *GoqueryProcessURL) parseSectionQueryType(ctx context.Context, scrapingQuery defn.ScrapeDataContentDefn, baseSelection *goquery.Selection) (string, []string) {
	log := util.GetGlobalLogger(ctx)
	var processedStrBuilder strings.Builder
	var relevantLinks []string

	if scrapingQuery.SectionType != nil {
		prefix := scrapingQuery.SectionType.Prefix
		suffix := scrapingQuery.SectionType.Suffix
		if strings.EqualFold(suffix, "") {
			suffix = "\n"
		}

		startElementSelection := baseSelection.Find(scrapingQuery.SectionType.StartSelector).First()
		// endElementSelection := baseSelection.Find(scrapingQuery.SectionType.EndSelector).First()

		finalSelection := startElementSelection.NextAll()

		if startElementSelection.Length() == 0 { // if start not found, don't do anything, if end not found, go through all of start's siblings
			finalSelection = baseSelection.Children()
			log.Println("no start selection found for scrapingquery", scrapingQuery.Name, "\nstarting with children of base selection")
		}

		titleElements := scrapingQuery.SectionType.Title
		if titleElements == nil {
			titleElements = []string{"h3"}
		}

		dataElements := scrapingQuery.SectionType.Data
		if dataElements == nil {
			dataElements = []string{"p"}
		}

		var endAchieved bool

		finalSelection.Each(func(i int, selection *goquery.Selection) {
			if endAchieved {
				return
			}
			if selection.Is(scrapingQuery.SectionType.EndSelector) {
				endAchieved = true
			}
			for _, title := range titleElements {
				if selection.Is(title) {
					text := selection.Text()
					if !strings.EqualFold(text, "") {
						relevantLinks = append(relevantLinks, util.GetLinksFromSelection(ctx, selection)...)
						processedStrBuilder.WriteString("## " + text + "\n")
						return
					}
				}
			}

			for _, content := range dataElements {
				if selection.Is(content) {
					relevantLinks = append(relevantLinks, util.GetLinksFromSelection(ctx, selection)...)
					if selection.Is("ol") {
						count, err := strconv.Atoi(selection.AttrOr("start", "1"))
						if err != nil {
							count = 1
						}

						selection.Children().Each(func(i int, itemSelector *goquery.Selection) {
							if itemSelector.Is("li") {
								processedStrBuilder.WriteString(fmt.Sprintf("%d. ", count) + itemSelector.Text() + suffix)
								count++
							} else {
								processedStrBuilder.WriteString(prefix + itemSelector.Text() + suffix)
							}
						})
					} else if selection.Is("ul") {
						selection.Children().Each(func(i int, itemSelector *goquery.Selection) {
							if itemSelector.Is("li") {
								processedStrBuilder.WriteString("* " + itemSelector.Text() + suffix)
							} else {
								processedStrBuilder.WriteString(prefix + itemSelector.Text() + suffix)
							}
						})
					} else {
						processedStrBuilder.WriteString(prefix + selection.Text() + suffix)
					}
					return
				}
			}
		})
	}

	return processedStrBuilder.String(), relevantLinks
}

func (processor *GoqueryProcessURL) scrapeImage(ctx context.Context, selection *goquery.Selection) []map[string]interface{} {
	log := util.GetGlobalLogger(ctx)
	var addedImages []string
	var imagesData []map[string]interface{}
	var imgIdx = 0

	taskId := processor.scrapeInfo["task-id"].(string)
	selection.Each(func(i int, imageSelector *goquery.Selection) {
		imageLink, found := imageSelector.Attr("src")
		if found {
			var imageName string = fmt.Sprintf("taskid_%s_image_%d", taskId, imgIdx)
			var requestImageLink string
			stringSplit := strings.Split(imageLink, "/")
			imageNameMeta := stringSplit[len(stringSplit)-1]
			imageExt := filepath.Ext(imageNameMeta)

			foundImageIdx := slices.Index(addedImages, imageLink)
			if foundImageIdx < 0 {
				addedImages = append(addedImages, imageLink)
				requestImageLink = strings.TrimLeft(imageLink, "/")
				if !strings.HasPrefix(requestImageLink, "http://") && !strings.HasPrefix(requestImageLink, "https://") {
					requestImageLink = "http://" + requestImageLink
				}

				imgResp, cerr := http.Get(requestImageLink)
				if cerr != nil {
					log.Println("could not retrieve image from link", imageLink, cerr)
					return
				}
				imgContent, err := io.ReadAll(imgResp.Body)
				if err != nil {
					fmt.Println("failed to copy img", imageName, "content")
					return
				}

				imagesData = append(imagesData, map[string]interface{}{
					"name":      imageName,
					"extension": imageExt,
					"content":   imgContent,
					"additionalinfo": map[string]interface{}{
						"url":         imageLink,
						"name":        imageNameMeta,
						"originalurl": requestImageLink,
					},
				})
				imgIdx++
			} else {
				// imageName = fmt.Sprintf("taskid_%s_image_%d", taskId, foundImageIdx)
			}
			// imageSelector.SetAttr(defn.ImageNameAttribute, imageName+imageExt)
		}
	})
	log.Println("successfully downloaded images for task", taskId)
	return imagesData
}

func (processor *GoqueryProcessURL) nestedScrape(ctx context.Context, links []string) (map[string]interface{}, *util.CustomError) {
	log := util.GetGlobalLogger(ctx)
	if processor.config.Depth > 0 {
		baseScrapeURL := processor.scrapeInfo["url"].(string)
		// log.Println("nested links for depth", processor.config.Depth, "are", links)

		parsedURL, err := url.Parse(baseScrapeURL)
		if err != nil { //this will never happen as initially this is checked
			cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeFailedToParseUrl, defn.ErrFailedToParseUrl, map[string]string{
				"error": err.Error(),
			})
			log.Println(cerr)
			return nil, cerr
		}

		baseUrl := fmt.Sprintf("%s://%s/", parsedURL.Scheme, parsedURL.Host)
		for _, hrefLink := range links {
			//if limit <= 0, scrape all links, else scrape the limit
			if processor.config.MaxLimit > 0 && len(processor.visitedLinks) >= processor.config.MaxLimit {
				break
			}

			followLink := ""
			parsedHrefLink, err := url.Parse(hrefLink)
			if err != nil {
				cerr := util.NewCustomError(ctx, "", err)
				log.Println("failed to parse provided url", cerr)
				continue
			}
			parsedHrefLink.Fragment = ""
			hrefLink = parsedHrefLink.String()

			if strings.EqualFold(parsedHrefLink.Scheme, "") && strings.EqualFold(parsedHrefLink.Host, "") {
				if strings.HasPrefix(hrefLink, "/") {
					followLink = baseUrl + strings.TrimPrefix(hrefLink, "/")
				} else {
					urlParts := strings.Split(strings.TrimSuffix(baseScrapeURL, "/"), "/")
					hrefParts := strings.Split(strings.TrimPrefix(hrefLink, "/"), "/")

					// get unique elements which are not present in url
					var uniqueHref []string
					for _, hrefPart := range hrefParts {
						if slices.Contains(urlParts, hrefPart) {
							uniqueHref = append(uniqueHref, hrefPart)
						}
					}

					// combine the unique elements with the original URL
					followLink = strings.Join(append(urlParts, uniqueHref...), "/")
				}
			} else if strings.HasPrefix(hrefLink, baseUrl) {
				followLink = hrefLink
			} else {
				log.Println("the link '", hrefLink, "' could not be opened because it is not in the same website as baselink")
				continue
			}

			followLink = strings.Trim(followLink, "/")
			if strings.EqualFold(followLink, "") {
				log.Println("skipping url", hrefLink, "as the link to follow is empty")
				continue
			} else if slices.Contains(processor.visitedLinks, followLink) {
				log.Println("skipping url", followLink, "as it has already been visited")
				continue
			}

			processor.visitedLinks = append(processor.visitedLinks, followLink)
			nestedConfig := processor.config
			nestedConfig.Depth = nestedConfig.Depth - 1

			currLevel := 0
			_, ok := processor.scrapeInfo["level"]
			if ok {
				currLevel = processor.scrapeInfo["level"].(int) + 1
			}
			scrapeInfo := map[string]interface{}{
				"url":                followLink,
				"visitedurls":        processor.visitedLinks,
				"level":              currLevel,
				"job-id":             processor.scrapeInfo["job-id"],
				"all_uploaded_files": processor.scrapeInfo["all_uploaded_files"],
			}

			var urlBasedScraper *UrlScraperService
			log.Println("starting to scrape child link", hrefLink, "at the url", followLink)
			if childScraper, cerr := urlBasedScraper.Init(ctx, nestedConfig, scrapeInfo); cerr != nil {
				log.Println("failed to initialise child link scraper for link", hrefLink, cerr)
				continue
			} else if returnedConfig, cerr := childScraper.SyncStart(ctx); cerr != nil {
				log.Println("failed to start child scrape job for link", hrefLink, cerr)
				continue
			} else {
				processor.visitedLinks = returnedConfig["visitedurls"].([]string)
				log.Println("childscrape response:", returnedConfig["response"])
			}
		}
		return map[string]interface{}{
			"visitedurls": processor.visitedLinks,
		}, nil
	}
	return nil, nil
}

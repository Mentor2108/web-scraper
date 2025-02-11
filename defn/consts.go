package defn

import "time"

const (
	ReadTimeout       = 15 * time.Second
	WriteTimeout      = 15 * time.Second
	ReadHeaderTimeout = 5 * time.Second

	ChromedpTimeout = 30 * time.Second
)

const (
	ContentTypeJSON         = "application/json"
	ContentTypePlainText    = "text/plain; charset=UTF-8"
	ContentTypeHTMLText     = "text/html; charset=utf-8"
	ContentTypeMarkdownText = "text/markdown"
	ContentTypeOctetStream  = "application/octet-stream"
	HTTPHeaderContentType   = "Content-Type"
)

const (
	ScrapePhaseLibraryChromedp = "chromedp"

	ProcessPhaseLibraryGoquery = "goquery"

	ScrapeQueryTypeSection = "section"
	ScrapeQueryTypeTable   = "table"
	ScrapeQueryTypeText    = "text"
)

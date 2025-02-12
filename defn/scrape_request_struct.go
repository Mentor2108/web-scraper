package defn

type ScrapeRequest struct {
	Url    string        `json:"url"`
	Config *ScrapeConfig `json:"config"`
}

type ScrapeConfig struct {
	Root              string                  `json:"root"`
	Depth             int                     `json:"depth"`
	MaxLimit          int                     `json:"maxlimit"`
	ContinueOnError   bool                    `json:"continueonerror"` //not using as of now
	ScrapePhase       *ScrapePhaseDefn        `json:"scrape_phase"`
	ProcessPhase      *ProcessPhaseDefn       `json:"process_phase"`
	ScrapeDataContent []ScrapeDataContentDefn `json:"content"`
	ExcludeElements   []string                `json:"exclude"`
	ScrapeImages      bool                    `json:"scrape_images"`
}

type ScrapePhaseDefn struct {
	Library string       `json:"library"`
	WaitFor *WaitForDefn `json:"wait_for"`
}

type ProcessPhaseDefn struct {
	Library string `json:"library"`
}

type ScrapeDataContentDefn struct {
	Name        string           `json:"name"`
	Type        string           `json:"type"`
	Selector    string           `json:"selector"`
	SectionType *SectionTypeDefn `json:"section"`
	TableType   *TableTypeDefn   `json:"table"`
	TextType    *TextTypeDefn    `json:"text"`
}

// WIP
type SectionTypeDefn struct {
	Prefix        string   `json:"prefix"`
	Suffix        string   `json:"suffix"`
	StartSelector string   `json:"start"`
	EndSelector   string   `json:"end"`
	Title         []string `json:"title"`
	Data          []string `json:"data"`
}

// WIP
type TableTypeDefn struct {
	Prefix           string         `json:"prefix"`
	Suffix           string         `json:"suffix"`
	Title            string         `json:"title"`
	ColumnsMap       *ColumnMapDefn `json:"column_map"`
	ColumnsNamesList []string       `json:"column_names"`
}

type ColumnMapDefn struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// WIP
type TextTypeDefn struct {
	Prefix string `json:"prefix"`
	Suffix string `json:"suffix"`
}

type WaitForDefn struct {
	Duration int    `json:"duration"`
	Selector string `json:"selector"`
}

func DefaultScrapeRequest() ScrapeRequest {
	return ScrapeRequest{
		Config: &ScrapeConfig{
			Root:              "body",
			Depth:             0,
			MaxLimit:          0,
			ContinueOnError:   true,
			ScrapeDataContent: []ScrapeDataContentDefn{},
			ExcludeElements: []string{
				"nav",
				"script",
				"noscript",
				"button",
				"style",
			},
			ScrapeImages: false,
		},
	}
}

package defn

type ScrapeRequest struct {
	Url    string        `json:"url"`
	Config *ScrapeConfig `json:"config"`
}

type ScrapeConfig struct {
	Root            string            `json:"root"`
	Depth           int               `json:"depth"`
	MaxLimit        int               `json:"maxlimit"`
	ContinueOnError bool              `json:"continueonerror"` //not using as of now
	ScrapePhase     *ScrapePhaseDefn  `json:"scrape_phase"`
	ProcessPhase    *ProcessPhaseDefn `json:"process_phase"`
	ScrapeData      []ScrapeDataDefn  `json:"scrape-data"`
}

type ScrapePhaseDefn struct {
	Library string       `json:"library"`
	WaitFor *WaitForDefn `json:"wait_for"`
}

type ProcessPhaseDefn struct {
	Library string `json:"library"`
}

type ScrapeDataDefn struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Selector string `json:"selector"`
}

type WaitForDefn struct {
	Duration int    `json:"duration"`
	Selector string `json:"selector"`
}

func DefaultScrapeRequest() ScrapeRequest {
	return ScrapeRequest{
		Config: &ScrapeConfig{
			Root:            "body",
			Depth:           1,
			MaxLimit:        1,
			ContinueOnError: true,
			ScrapeData:      []ScrapeDataDefn{},
		},
	}
}

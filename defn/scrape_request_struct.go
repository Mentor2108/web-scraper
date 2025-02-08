package defn

type ScrapeRequest struct {
	Url    string        `json:"url"`
	Config *ScrapeConfig `json:"config"`
}

type ScrapeConfig struct {
	Root            string            `json:"root"`
	Depth           int               `json:"depth"`
	MaxLimit        int               `json:"maxlimit"`
	ContinueOnError bool              `json:"continueonerror"`
	ScrapePhase     *ScrapePhaseDefn  `json:"scrape-phase"`
	ProcessPhase    *ProcessPhaseDefn `json:"process-phase"`
	ScrapeData      []ScrapeDataDefn  `json:"scrape-data"`
}

type ScrapePhaseDefn struct {
	Library string `json:"library"`
}

type ProcessPhaseDefn struct {
	Library string `json:"library"`
}

type ScrapeDataDefn struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Selector string `json:"selector"`
}

package defn

type ScrapeJob struct {
	ID       string                 `json:"id"`
	URL      string                 `json:"url"`
	Depth    int                    `json:"depth"`
	Maxlimit int                    `json:"maxlimit"`
	Response map[string]interface{} `json:"response"` // JSONB Field
}

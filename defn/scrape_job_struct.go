package defn

type ScrapeJob struct {
	ID       string                 `json:"id"`
	Depth    int                    `json:"depth"`
	Maxlimit int                    `json:"maxlimit"`
	Response map[string]interface{} `json:"response"` // JSONB Field
}

package defn

type ScrapeTask struct {
	ID       string                 `json:"id"`
	JobId    string                 `json:"job_id"`
	URL      string                 `json:"url"`
	Depth    int                    `json:"depth"`
	Maxlimit int                    `json:"maxlimit"`
	Level    int                    `json:"level"`
	Response map[string]interface{} `json:"response"` // JSONB Field
}

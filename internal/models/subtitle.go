package models

type LineFormat struct {
	Prefix string
	Suffix string
}

type SubtitleBlock struct {
	ID         string
	Time       string
	Text       string
	Formatting []LineFormat
}

type JobEstimate struct {
	TotalLines       int     `json:"total_lines"`
	CachedLines      int     `json:"cached_lines"`
	LinesToTranslate int     `json:"lines_to_translate"`
	TotalBatches     int     `json:"total_batches"`
	EstimatedTokens  int     `json:"estimated_tokens"`
	EstimatedCostUSD float64 `json:"estimated_cost_usd"`
}

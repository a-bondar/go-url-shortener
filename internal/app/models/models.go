package models

type HandleShortenRequest struct {
	URL string `json:"url"`
}

type HandleShortenResponse struct {
	Result string `json:"result"`
}

type OriginalURLCorrelation struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ShortURLCorrelation struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type HandleShortenBatchRequest []OriginalURLCorrelation

type HandleShortenBatchResponse []ShortURLCorrelation

type Data struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
	Deleted     bool   `json:"deleted"`
}

type URLsPair struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type HandleUserURLsResponse []URLsPair

package models

type HandleShortenRequest struct {
	URL string `json:"url"`
}

type HandleShortenResponse struct {
	Result string `json:"result"`
}

type HandleShortenBatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type HandleShortenBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type Data struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

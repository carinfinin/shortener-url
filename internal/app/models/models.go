package models

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type Line struct {
	URL string `json:"url"`
	ID  string `json:"id"`
}

type RequestBatch struct {
	ID      string `json:"correlation_id"`
	LongURL string `json:"original_url"`
}

type ResponseBatch struct {
	ID       string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

type AuthLine struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
	IsDeleted   string `json:"is_deleted"`
}

type UserURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

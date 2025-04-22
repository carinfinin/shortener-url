package models

// Request структура для отбработки json реквеста.
type Request struct {
	URL string `json:"url"`
}

// Response структура для отбработки json ответа.
type Response struct {
	Result string `json:"result"`
}

//type Line struct {
//	URL string `json:"url"`
//	ID  string `json:"id"`
//}

// RequestBatch структупа для обработки запроса json.
type RequestBatch struct {
	ID      string `json:"correlation_id"`
	LongURL string `json:"original_url"`
}

// RequestBatch структупа для обработки ответа json.
type ResponseBatch struct {
	ID       string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

// AuthLine записывается в бд.
type AuthLine struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
	IsDeleted   bool   `json:"is_deleted"`
}

// UserURL структыра для обработки ответа получения урл пользователя.
type UserURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// DeleteURLUser структыра для обработки удаления урлов пользователя.
type DeleteURLUser struct {
	Data   string
	USerID string
}

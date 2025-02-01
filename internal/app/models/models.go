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

package storage

type Repositories interface {
	AddUrl(url string) string
	GetUrl(xmlID string) (string, error)
}

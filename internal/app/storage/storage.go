package storage

type Repositories interface {
	AddURL(url string) string
	GetURL(xmlID string) (string, error)
}

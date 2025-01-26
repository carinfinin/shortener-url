package storage

type Repositories interface {
	AddURL(url string) (string, error)
	GetURL(xmlID string) (string, error)
}

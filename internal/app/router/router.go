package router

import (
	"github.com/carinfinin/shortener-url/internal/app/storage"
	"io"
	"net/http"
	"strings"
)

type Router struct {
	Handle *http.ServeMux
	Store  storage.Repositories
}

func ConfigureRouter(s storage.Repositories) *Router {

	r := Router{
		Handle: http.NewServeMux(),
		Store:  s,
	}

	r.Handle.HandleFunc(http.MethodPost+" /", CreateURL(r))
	r.Handle.HandleFunc(http.MethodGet+" /{id}", GetURL(r))

	return &r
}

func CreateURL(r Router) http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		defer req.Body.Close()

		url := strings.TrimSpace(string(body))

		xmlID := r.Store.AddURL(url)
		newURL := "http://localhost:8080/" + xmlID

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(newURL))
	}
}

func GetURL(r Router) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		path := strings.TrimPrefix(req.URL.Path, "/")
		path = strings.TrimSuffix(path, "/")
		parts := strings.Split(path, "/")

		if len(parts) != 1 {
			http.Error(res, "Not found", http.StatusBadRequest)
			return
		}
		xmlID := parts[0]

		if xmlID == "" {
			http.Error(res, "Not found", http.StatusBadRequest)
			return
		}

		url, err := r.Store.GetURL(xmlID)
		if err != nil {
			http.NotFound(res, req)
			return
		}

		//log.Printf("Retrieved URL: %s", url)

		if url == "" {
			http.NotFound(res, req)
			return
		}

		//res.Header().Set("Location", url)
		//res.WriteHeader(http.StatusTemporaryRedirect)
		http.Redirect(res, req, url, http.StatusTemporaryRedirect)

	}
}

package router

import (
	middleware2 "github.com/carinfinin/shortener-url/internal/app/middleware"
	"github.com/carinfinin/shortener-url/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strings"
)

type Router struct {
	Handle *chi.Mux
	Store  storage.Repositories
	URL    string
}

func ConfigureRouter(s storage.Repositories, url string) *Router {

	r := Router{
		Handle: chi.NewRouter(),
		Store:  s,
		URL:    url,
	}

	r.Handle.Use(middleware2.Logging)

	r.Handle.Post("/", CreateURL(r))
	r.Handle.Get("/{id}", GetURL(r))

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
		newURL := r.URL + "/" + xmlID

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(newURL))
	}
}

func GetURL(r Router) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		path := strings.Trim(req.URL.Path, "/")
		parts := strings.Split(path, "/")
		id := parts[0]

		if id == "" {
			http.Error(res, "Not found", http.StatusBadRequest)
			return
		}

		url, err := r.Store.GetURL(id)
		if err != nil {
			http.NotFound(res, req)
			return
		}

		if url == "" {
			http.NotFound(res, req)
			return
		}

		http.Redirect(res, req, url, http.StatusTemporaryRedirect)
	}
}

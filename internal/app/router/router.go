package router

import (
	"encoding/json"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	middleware2 "github.com/carinfinin/shortener-url/internal/app/middleware"
	"github.com/carinfinin/shortener-url/internal/app/models"
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

	r.Handle.Use(middleware2.RequestLogger)
	r.Handle.Use(middleware2.ResponseLogger)

	r.Handle.Post("/api/shorten", JSONHandle(r))
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

func JSONHandle(r Router) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost {

			logger.Log.Info("start handle JSON")

			var req models.Request
			decoder := json.NewDecoder(request.Body)
			err := decoder.Decode(&req)
			if err != nil {
				logger.Log.Error("Decode error", err)
				http.Error(writer, "bad request", http.StatusBadRequest)
				return
			}
			req.Url = strings.TrimSpace(req.Url)

			xmlID := r.Store.AddURL(req.Url)

			var res models.Response

			res.Result = r.URL + "/" + xmlID

			encoder := json.NewEncoder(writer)

			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusCreated)

			if err := encoder.Encode(res); err != nil {
				logger.Log.Error("Encode error", err)
				http.Error(writer, "bad request", http.StatusBadRequest)
				return
			}

		} else {
			http.NotFound(writer, request)
			return
		}

	}
}

package router

import (
	"encoding/json"
	"errors"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	middleware2 "github.com/carinfinin/shortener-url/internal/app/middleware"
	"github.com/carinfinin/shortener-url/internal/app/models"
	"github.com/carinfinin/shortener-url/internal/app/storage"
	"github.com/carinfinin/shortener-url/internal/app/storage/storepg"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strings"
)

type Router struct {
	Handle *chi.Mux
	Store  storage.Repositories
	URL    string
	Config *config.Config
}

func ConfigureRouter(s storage.Repositories, config *config.Config) *Router {

	r := Router{
		Handle: chi.NewRouter(),
		Store:  s,
		URL:    config.URL,
		Config: config,
	}
	r.Handle.Use(middleware2.CompressGzipWriter)
	r.Handle.Use(middleware2.CompressGzipReader)
	r.Handle.Use(middleware2.RequestLogger)
	r.Handle.Use(middleware2.ResponseLogger)

	r.Handle.Post("/api/shorten", JSONHandle(r))
	r.Handle.Post("/api/shorten/batch", JSONHandleBatch(r))
	r.Handle.Post("/", CreateURL(r))
	r.Handle.Get("/ping", PingDB(r))
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

		xmlID, err := r.Store.AddURL(url)

		newURL := r.URL + "/" + xmlID
		res.Header().Set("Content-Type", "text/plain")

		if err != nil {
			if errors.Is(err, storage.ErrDouble) {
				res.WriteHeader(http.StatusConflict)
				res.Write([]byte(newURL))
				return
			}
			logger.Log.Error("CreateURL", err)
			http.Error(res, "CreateURL", http.StatusInternalServerError)
			return
		}

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

		logger.Log.Info("start handle JSON")

		var req models.Request
		decoder := json.NewDecoder(request.Body)
		err := decoder.Decode(&req)
		if err != nil {
			logger.Log.Error("Decode error", err)
			http.Error(writer, "bad request", http.StatusBadRequest)
			return
		}
		req.URL = strings.TrimSpace(req.URL)

		xmlID, err := r.Store.AddURL(req.URL)
		if err != nil && len(xmlID) == 0 {

			logger.Log.Error("JSONHandle", err)
			http.Error(writer, "error add url", http.StatusInternalServerError)
			return
		}

		var res models.Response

		res.Result = r.URL + "/" + xmlID

		encoder := json.NewEncoder(writer)

		writer.Header().Set("Content-Type", "application/json")

		if err != nil && errors.Is(err, storage.ErrDouble) {
			writer.WriteHeader(http.StatusConflict)
		} else {
			writer.WriteHeader(http.StatusCreated)
		}

		if err := encoder.Encode(res); err != nil {
			logger.Log.Error("Encode error", err)
			http.Error(writer, "bad request", http.StatusBadRequest)
			return
		}

	}
}

func JSONHandleBatch(r Router) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		logger.Log.Info("start handle JSONHandleBatch")

		var data = make([]models.RequestBatch, 0)

		decoder := json.NewDecoder(request.Body)

		err := decoder.Decode(&data)
		if err != nil {
			logger.Log.Error("Decode error", err)
			http.Error(writer, "bad request", http.StatusBadRequest)
			return
		}

		result, err := r.Store.AddURLBatch(data)
		if err != nil {
			logger.Log.Error("JSONHandleBatch AddURLBatch error", err)
			http.Error(writer, "bad request", http.StatusBadRequest)
			return
		}

		encoder := json.NewEncoder(writer)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)

		if err := encoder.Encode(result); err != nil {
			logger.Log.Error("Encode error", err)
			http.Error(writer, "bad request", http.StatusBadRequest)
			return
		}

	}
}

func PingDB(r Router) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		logger.Log.Info("PingDB handler start")

		err := storepg.Ping(r.Config.DBPath)
		if err != nil {
			logger.Log.Info("PingDB handler error: ", err)
			writer.WriteHeader(http.StatusInternalServerError)
		}
		logger.Log.Info("PingDB handler status OK")

		writer.WriteHeader(http.StatusOK)

	}
}

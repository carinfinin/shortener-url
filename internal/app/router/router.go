package router

import (
	"encoding/json"
	"errors"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	middleware2 "github.com/carinfinin/shortener-url/internal/app/middleware"
	"github.com/carinfinin/shortener-url/internal/app/models"
	"github.com/carinfinin/shortener-url/internal/app/service"
	"github.com/carinfinin/shortener-url/internal/app/storage"
	"github.com/carinfinin/shortener-url/internal/app/storage/storepg"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strings"
)

type Router struct {
	Handle  *chi.Mux
	URL     string
	Config  *config.Config
	Service *service.Service
}

func ConfigureRouter(s *service.Service, config *config.Config) *Router {

	r := Router{
		Handle:  chi.NewRouter(),
		URL:     config.URL,
		Config:  config,
		Service: s,
	}
	r.Handle.Use(middleware2.CompressGzipWriter)
	r.Handle.Use(middleware2.CompressGzipReader)
	r.Handle.Use(middleware2.RequestLogger)
	r.Handle.Use(middleware2.ResponseLogger)
	r.Handle.Use(middleware2.AuthMiddleWare)

	r.Handle.Post("/api/shorten", JSONHandle(r))
	r.Handle.Post("/api/shorten/batch", JSONHandleBatch(r))
	r.Handle.Post("/", CreateURL(r))
	r.Handle.Get("/ping", PingDB(r))
	r.Handle.Get("/{id}", GetURL(r))
	r.Handle.Get("/api/user/urls", GetUserURLs(r))
	r.Handle.Delete("/api/user/urls", DeleteUserURLs(r))

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

		xmlID, err := r.Service.CreateURL(req.Context(), url)

		newURL := r.URL + "/" + xmlID
		res.Header().Set("Content-Type", "text/plain")

		if err != nil {
			if errors.Is(err, storage.ErrDouble) {
				res.WriteHeader(http.StatusConflict)
				res.Write([]byte(newURL))
				return
			}
			logger.Log.Error("CreateURL error: ", err)
			http.Error(res, "CreateURL error", http.StatusInternalServerError)
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

		url, err := r.Service.GetURL(req.Context(), id)
		if err != nil {
			if errors.Is(err, storage.ErrDeleteURL) {
				http.Error(res, "URL is deleted", http.StatusGone)
				return
			}
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

		xmlID, err := r.Service.JSONHandle(request.Context(), req.URL)
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

		result, err := r.Service.JSONHandleBatch(request.Context(), data)
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

func GetUserURLs(r Router) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		logger.Log.Info("GetUserURLs handler start")

		encoder := json.NewEncoder(writer)

		data, err := r.Service.GetUserURLs(request.Context())
		if err != nil {
			logger.Log.Info("GetUserURLs handler error: ", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(data) == 0 {
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		if err := encoder.Encode(data); err != nil {
			logger.Log.Error("Encode error", err)
			http.Error(writer, "server error", http.StatusBadRequest)
			return
		}
	}
}

func DeleteUserURLs(r Router) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		logger.Log.Info("DeleteUserURLs handler start")

		var data = make([]string, 0)
		decoder := json.NewDecoder(request.Body)

		err := decoder.Decode(&data)
		if err != nil {
			logger.Log.Error("Decode error", err)
			http.Error(writer, "bad request", http.StatusBadRequest)
			return
		}
		r.Service.DeleteUserURLs(request.Context(), data)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusAccepted)
	}
}

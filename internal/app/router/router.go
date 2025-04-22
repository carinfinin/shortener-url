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

// Router represents the HTTP router application
type Router struct {
	Handle  *chi.Mux         // chi router for handling HTTP requests
	URL     string           // router base URL
	Config  *config.Config   // application configuration
	Service service.IService // service layer with business logic
}

// ConfigureRouter constructor for type Router accepts *service.Service *config.Config.
func ConfigureRouter(s service.IService, config *config.Config) *Router {

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

	r.Handle.Post("/api/shorten", r.JSONHandle)
	r.Handle.Post("/api/shorten/batch", r.JSONHandleBatch)
	r.Handle.Post("/", r.CreateURL)
	r.Handle.Get("/ping", r.PingDB)
	r.Handle.Get("/{id}", r.GetURL)
	r.Handle.Get("/api/user/urls", r.GetUserURLs)
	r.Handle.Delete("/api/user/urls", r.DeleteUserURLs)

	return &r
}

// CreateURL creates a new url.
//
// Accepts POST request in text/plain:
// https://practicum.yandex.ru/learn/go-advanced/courses/
//
// Possible response codes:
// - 201 Created - record successfully created
// - 400 Bad Request - invalid request format
// - 409 Conflict - URL already exists
// - 500 Internal Server Error - server error
func (r *Router) CreateURL(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
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

// GetURL receiving url.
//
// Possible response codes:
// - 400 Bad Request - invalid request format
// - 410 url deleted
// - 307 redirect
func (r *Router) GetURL(res http.ResponseWriter, req *http.Request) {

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

// JSONHandle processes a JSON request to create a new record.
//
// Accepts a POST request with a JSON body in the format:
// {
// "url": "string" // URL to process
// }
//
// Possible response codes:
// - 201 Created - the record was successfully created
// - 400 Bad Request - invalid request format
// - 409 Conflict - the URL already exists
// - 500 Internal Server Error - server error
//
// If successful, returns a JSON response:
// {
// "result": "string" // Generated URL with XML ID
// }
func (r *Router) JSONHandle(writer http.ResponseWriter, request *http.Request) {
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

// JSONHandleBatch processes a JSON request to create multiple new records.
//
// Accepts a POST request with a JSON body in the format:
// [{
// "correlation_id": "string"
// "original_url": "string"
// }]
//
// Possible response codes:
// - 201 Created - the record was successfully created
// - 400 Bad Request - invalid request format
// - 500 Internal Server Error - server error
func (r *Router) JSONHandleBatch(writer http.ResponseWriter, request *http.Request) {
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

func (r *Router) PingDB(writer http.ResponseWriter, request *http.Request) {
	logger.Log.Info("PingDB handler start")

	err := storepg.Ping(r.Config.DBPath)
	if err != nil {
		logger.Log.Info("PingDB handler error: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
	logger.Log.Info("PingDB handler status OK")

	writer.WriteHeader(http.StatusOK)

}

// GetUserURLs receives user generated URLs.
//
// Possible response codes:
// - 200 Success
// - 204 No content
// - 400 Bad Request - invalid request format
// - 500 Internal Server Error - server error
//
// If successful, returns a JSON response:
// [{
// "short_url": "string"
// "original_url": "string"
// }]
func (r *Router) GetUserURLs(writer http.ResponseWriter, request *http.Request) {
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

// DeleteUserURLs processes a JSON request to delete records.
//
// Accepts a POST request with a JSON body in the format:
// ["short_url", "short_url_2"]
//
// Possible response codes:
// - 202 Accepted
// - 400 Bad Request - invalid request format
func (r *Router) DeleteUserURLs(writer http.ResponseWriter, request *http.Request) {
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

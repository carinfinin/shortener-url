package router

import (
	"github.com/carinfinin/shortener-url/internal/app/storage"
	"io"
	"net/http"
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

	r.Handle.HandleFunc(http.MethodPost+" /", r.createURL)
	r.Handle.HandleFunc(http.MethodGet+" /{id}/", r.getURL)

	return &r
}

func (r *Router) createURL(res http.ResponseWriter, req *http.Request) {

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	xmlID := r.Store.AddURL(string(body))
	newURL := "http://localhost:8080/" + xmlID

	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(newURL))
}

func (r *Router) getURL(res http.ResponseWriter, req *http.Request) {

	xmlID := req.PathValue("id")

	if xmlID == "" {
		//http.NotFound(res, req)
		http.Error(res, "Not found", http.StatusBadRequest)
	} else {

		//res.Header().Set()

		url, err := r.Store.GetURL(xmlID)
		if err != nil {
			http.NotFound(res, req)
		}

		res.Header().Set("Location", url)
		res.WriteHeader(http.StatusTemporaryRedirect)

		//http.Redirect(res, req, url, http.StatusMovedPermanently)
	}

}

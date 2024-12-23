package router

import (
	"fmt"
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

	r.Handle.HandleFunc(http.MethodPost+" /", r.createUrl)
	r.Handle.HandleFunc(http.MethodGet+" /{id}/", r.getUrl)

	return &r
}

func (r *Router) createUrl(res http.ResponseWriter, req *http.Request) {

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	xmlID := r.Store.AddUrl(string(body))
	fmt.Fprint(res, xmlID)
}

func (r *Router) getUrl(res http.ResponseWriter, req *http.Request) {

	xmlId := req.PathValue("id")

	if xmlId == "" {
		//http.NotFound(res, req)
		http.Error(res, "Not found", http.StatusBadRequest)
	} else {

		//res.Header().Set()

		url, err := r.Store.GetUrl(xmlId)
		if err != nil {
			fmt.Println(err)
			http.NotFound(res, req)
		}

		http.Redirect(res, req, url, http.StatusMovedPermanently)
		//res.Write([]byte(url))
	}

}

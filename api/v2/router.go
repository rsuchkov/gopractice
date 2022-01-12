package v2

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/rsuchkov/gopractice/service/serverstats"
)

func NewRouter(svc *serverstats.Processor) chi.Router {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/update/", func(rw http.ResponseWriter, r *http.Request) {
		UpdateMetricHandler(svc, rw, r)
	})

	r.Post("/update", func(rw http.ResponseWriter, r *http.Request) {
		UpdateMetricHandler(svc, rw, r)
	}) // this looks quite stupid, but I couldnt find a way to do it correct
	r.Post("/value/", func(rw http.ResponseWriter, r *http.Request) {
		GetMetricHandler(svc, rw, r)
	})
	return r
}

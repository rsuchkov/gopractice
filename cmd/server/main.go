package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/rsuchkov/gopractice/model"
	"github.com/rsuchkov/gopractice/service/serverstats"
	"github.com/rsuchkov/gopractice/storage/memory"
)

const (
	addr = "127.0.0.1:8080"
)

func UpdateMetricHandler(svc *serverstats.Processor, w http.ResponseWriter, r *http.Request) {

	mtype, name := model.MetricType(chi.URLParam(r, "mtype")), chi.URLParam(r, "mname")
	value, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
	if err != nil {
		http.Error(w, "Wrong metric value", http.StatusBadRequest)
		return
	}
	er := svc.SaveMetric(name, mtype, value)
	if er != nil {
		http.Error(w, "Wrong metric value", http.StatusNotImplemented)
		return
	}
	w.Write([]byte("ok"))
}

func GetMetricHandler(svc *serverstats.Processor, w http.ResponseWriter, r *http.Request) {
	mtype, name := model.MetricType(chi.URLParam(r, "mtype")), chi.URLParam(r, "mname")
	v, err := svc.GetMetric(name, mtype)
	if err != nil {
		http.Error(w, "Doesn't exist", http.StatusNotFound)
		return
	}
	w.Write([]byte(fmt.Sprintf("%f", v)))
}

func GetMetricsHandler(svc *serverstats.Processor, w http.ResponseWriter, r *http.Request) {

	metrics, _ := svc.GetMetrics()
	v, err := json.Marshal(metrics)
	if err != nil {
		http.Error(w, "Wrong data", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(string(v)))
}

func NewRouter(svc *serverstats.Processor) chi.Router {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/update/{mtype}/{mname}/{value}", func(rw http.ResponseWriter, r *http.Request) {
		UpdateMetricHandler(svc, rw, r)
	})
	r.Get("/value/{mtype}/{mname}", func(rw http.ResponseWriter, r *http.Request) {
		GetMetricHandler(svc, rw, r)
	})
	r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		GetMetricsHandler(svc, rw, r)
	})
	return r
}

func main() {
	st, err := memory.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	svc, err := serverstats.New(serverstats.WithStatsStorage(st))
	if err != nil {
		fmt.Println(err)
		return
	}
	r := NewRouter(svc)
	log.Fatal(http.ListenAndServe(addr, r))
}

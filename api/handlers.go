package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/rsuchkov/gopractice/model"
	"github.com/rsuchkov/gopractice/service/serverstats"
)

func UpdateMetricHandler(svc *serverstats.Processor, w http.ResponseWriter, r *http.Request) {

	mtype, name := model.MType(chi.URLParam(r, "mtype")), chi.URLParam(r, "mname")
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
	mtype, name := model.MType(chi.URLParam(r, "mtype")), chi.URLParam(r, "mname")
	v, err := svc.GetMetric(name, mtype)
	if err != nil {
		http.Error(w, "Doesn't exist", http.StatusNotFound)
		return
	}
	if mtype == model.MetricTypeGauge {
		w.Write([]byte(fmt.Sprint(v)))
	} else {
		w.Write([]byte(fmt.Sprint(int(v))))
	}

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

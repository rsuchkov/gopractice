package v2

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rsuchkov/gopractice/model"
	"github.com/rsuchkov/gopractice/service/serverstats"
)

func UpdateMetricHandler(svc *serverstats.Processor, w http.ResponseWriter, r *http.Request) {

	var m model.Metric
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if m.Value == nil && m.MType == model.MetricTypeGauge {
		http.Error(w, "Field Value has to be set", http.StatusBadRequest)
		return
	} else if m.Value == nil && m.MType == model.MetricTypeCounter {
		*m.Value = 0
	}

	ret, er := svc.SaveMetric(m)
	if er != nil {
		http.Error(w, "Wrong metric value", http.StatusNotImplemented)
		return
	}
	v, err := json.Marshal(ret)
	if err != nil {
		http.Error(w, "Wrong data", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(string(v)))
}

func GetMetricHandler(svc *serverstats.Processor, w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	var m model.Metric
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	ret, err := svc.GetMetric(m.ID, m.MType)
	if err != nil {
		http.Error(w, "Doesn't exist", http.StatusNotFound)
		return
	}
	v, err := json.Marshal(ret)
	if err != nil {
		http.Error(w, "Wrong data", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(string(v)))

}

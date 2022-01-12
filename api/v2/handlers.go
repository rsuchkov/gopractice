package v2

import (
	"encoding/json"
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

	ret, er := svc.SaveMetric_v2(m)
	if er != nil {
		http.Error(w, "Wrong metric value", http.StatusNotImplemented)
		return
	}
	v, err := json.Marshal(ret)
	if err != nil {
		http.Error(w, "Wrong data", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(string(v)))
}

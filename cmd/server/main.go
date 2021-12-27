package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/rsuchkov/gopractice/model"
	"github.com/rsuchkov/gopractice/service/serverstats"
	"github.com/rsuchkov/gopractice/storage/memory"
)

const (
	addr = "127.0.0.1:8080"
)

var metricRoute = regexp.MustCompile(`^/update/([a-z]+)/([a-zA-Z]+)/([a-zA-Z0-9\.]+)/?$`)

func Handler(svc *serverstats.Processor, w http.ResponseWriter, r *http.Request) {
	data := metricRoute.FindStringSubmatch(r.RequestURI)
	if len(data) != 4 {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Not allowed", http.StatusMethodNotAllowed)
		return
	}
	mtype, name := model.MetricType(data[1]), data[2]
	value, err := strconv.ParseFloat(data[3], 64)
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Handler(svc, w, r)
	})
	server := &http.Server{
		Addr: addr,
	}
	server.ListenAndServe()
}

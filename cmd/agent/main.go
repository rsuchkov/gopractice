package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/rsuchkov/gopractice/provider/mothership"
	"github.com/rsuchkov/gopractice/service/agentstats"
	"github.com/rsuchkov/gopractice/storage/memory"
)

const (
	server = "http://127.0.0.1:8080/"
)

type Monitor struct {
	Alloc         uint64
	BuckHashSys   uint64
	Frees         uint64
	GCCPUFraction float64
	GCSys         uint64
	HeapAlloc     uint64
	HeapIdle      uint64
	HeapInuse     uint64
	HeapObjects   uint64
	HeapReleased  uint64
	HeapSys       uint64
	LastGC        uint64
	Lookups       uint64
	MCacheInuse   uint64
	MCacheSys     uint64
	MSpanInuse    uint64
	MSpanSys      uint64
	Mallocs       uint64
	NextGC        uint64
	NumForcedGC   uint32
	NumGC         uint64
	OtherSys      uint64
	PauseTotalNs  uint64
	StackInuse    uint64
	StackSys      uint64
	Sys           uint64

	PollCount   uint32
	RandomValue int
}

func SendMetrics(duration int, m *Monitor) {
	var interval = time.Duration(duration) * time.Second

	for {
		<-time.After(interval)
		b, _ := json.Marshal(m)
		fmt.Println(string(b))

		client := &http.Client{}
		data := url.Values{}
		request, err := http.NewRequest(http.MethodPost, server, bytes.NewBufferString(data.Encode()))
		if err != nil {
			fmt.Println(err)
			continue
		}
		request.Header.Add("application-type", "text/plain")
		request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
		response, err := client.Do(request)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(string(body))
	}
}

func RuntimeMonitor(duration int, m *Monitor) {
	var rtm runtime.MemStats
	var interval = time.Duration(duration) * time.Second

	for {
		<-time.After(interval)

		runtime.ReadMemStats(&rtm)

		m.Alloc = rtm.Alloc
		m.BuckHashSys = rtm.BuckHashSys
		m.Frees = rtm.Frees
		m.GCCPUFraction = rtm.GCCPUFraction
		m.GCSys = rtm.GCSys
		m.HeapAlloc = rtm.HeapAlloc
		m.HeapIdle = rtm.HeapIdle
		m.HeapInuse = rtm.HeapInuse
		m.HeapObjects = rtm.HeapObjects
		m.HeapReleased = rtm.HeapReleased
		m.HeapSys = rtm.HeapSys
		m.LastGC = rtm.LastGC
		m.Lookups = rtm.Lookups
		m.MCacheInuse = rtm.MCacheInuse
		m.MCacheSys = rtm.MCacheSys
		m.MSpanInuse = rtm.MSpanInuse
		m.MSpanSys = rtm.MSpanSys
		m.Mallocs = rtm.Mallocs
		m.NextGC = rtm.NextGC
		m.NumForcedGC = rtm.NumForcedGC
		m.NumGC = rtm.NextGC
		m.OtherSys = rtm.OtherSys
		m.PauseTotalNs = rtm.PauseTotalNs
		m.StackInuse = rtm.StackInuse
		m.StackSys = rtm.StackSys
		m.Sys = rtm.Sys

		m.RandomValue = rand.Intn(100)
		m.PollCount += 1
	}
}

func main() {
	// var m Monitor
	// f := func() { SendMetrics(3, &m) }
	// time.AfterFunc(0, f)
	// RuntimeMonitor(2, &m)
	st, err := memory.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	svc, err := agentstats.New(agentstats.WithStatsStorage(st))
	if err != nil {
		fmt.Println(err)
		return
	}

	tickPoll := time.NewTicker(2 * time.Second)
	defer tickPoll.Stop()

	tickReport := time.NewTicker(10 * time.Second)
	defer tickReport.Stop()

	done := make(chan struct{})
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for sig := range sigCh {
			log.Println("Recieved sig:", sig)
			tickPoll.Stop()
			//tickReport.Stop()
			close(done)
			return
		}

	}()
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("Agent stopped")
				return
			case <-tickPoll.C:
				svc.CollectMetrics()
			}
		}
	}()

	// Report every 10s
	client := mothership.New(server)
	for {
		select {
		case <-done:
			fmt.Println("Aggent stopped")
			return
		case <-tickReport.C:
			for _, v := range st.GetMetrics() {
				select {
				case <-done:
					return
				default:
					time.Sleep(500 * time.Microsecond)
					resp, err := client.SendMetric(v)

					if err != nil {
						log.Println(err)
						log.Println("Failed to send", v.Name)

					}
					if resp != 0 {
						log.Println(resp, v.Name)
					}
				}
			}

		}
	}
}

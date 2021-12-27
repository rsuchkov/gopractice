package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rsuchkov/gopractice/provider/mothership"
	"github.com/rsuchkov/gopractice/service/agentstats"
	"github.com/rsuchkov/gopractice/storage/memory"
)

const (
	server = "http://127.0.0.1:8080/"
)

func main() {
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
			tickReport.Stop()
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

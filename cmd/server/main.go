package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rsuchkov/gopractice/api"
	"github.com/rsuchkov/gopractice/api/v2"
	"github.com/rsuchkov/gopractice/service/serverstats"
	"github.com/rsuchkov/gopractice/storage/memory"
)

const (
	addr = "127.0.0.1:8080"
)

func main() {
	st, err := memory.New()
	if err != nil {
		log.Fatal(err)
		return
	}
	svc, err := serverstats.New(serverstats.WithStatsStorage(st))
	if err != nil {
		log.Fatal(err)
		return
	}
	r := api.NewRouter(svc)
	r.Mount("/", v2.NewRouter(svc))

	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		for sig := range sigCh {
			log.Println("Recieved sig:", sig)
			fmt.Println("Dying...")
			server.Shutdown(context.Background())
		}

	}()

	log.Println("Starting on:", addr)
	log.Fatal(server.ListenAndServe())
}

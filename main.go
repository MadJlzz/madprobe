package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/madjlzz/madprobe/controller"
	"github.com/madjlzz/madprobe/internal/alerter"
	"github.com/madjlzz/madprobe/internal/prober"
	"github.com/madjlzz/madprobe/util"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	configuration := util.NewServerConfiguration()

	client := http.DefaultClient
	if len(configuration.CaCertificate) > 0 {
		client = util.HttpsClient(configuration.CaCertificate)
	}

	// Event Bus channel to let services communicate.
	alertBus := make(chan prober.Probe)
	probeService := prober.NewProbeService(client, alertBus)
	probeController := controller.NewProbeController(probeService)

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/probe/create", probeController.Create).
		Methods(http.MethodPost)
	r.HandleFunc("/api/v1/probe/{name}", probeController.Read).
		Methods(http.MethodGet)
	r.HandleFunc("/api/v1/probe", probeController.ReadAll).
		Methods(http.MethodGet)
	r.HandleFunc("/api/v1/probe/{name}", probeController.Update).
		Methods(http.MethodPut)
	r.HandleFunc("/api/v1/probe/{name}", probeController.Delete).
		Methods(http.MethodDelete)

	srv := &http.Server{
		Addr: "0.0.0.0:" + configuration.Port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if len(configuration.ServerCertificate) == 0 || len(configuration.ServerKey) == 0 {
			log.Printf("Starting HTTP server on port %s...\n", configuration.Port)
			if err := srv.ListenAndServe(); err != nil {
				log.Println(err)
			}
		} else {
			log.Printf("Starting HTTPs server on port %s...\n", configuration.Port)
			if err := srv.ListenAndServeTLS(configuration.ServerCertificate, configuration.ServerKey); err != nil {
				log.Println(err)
			}
		}
	}()

	// Alerter start at boot time.
	al, err := alerter.NewService(alertBus)
	if err != nil {
		fmt.Printf("[WARNING] alerter module wasn't able to start. got: %v\n", err)
	}
	al.Run()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), configuration.Wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	_ = srv.Shutdown(ctx)
	_ = al.Close()
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

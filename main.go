package main

import (
	"context"
	"flag"
	"github.com/gorilla/mux"
	"github.com/madjlzz/madprobe/controller"
	"github.com/madjlzz/madprobe/internal"
	"github.com/madjlzz/madprobe/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	wait = flag.Duration("graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	port = flag.String("port", "3000", "the port for which the server will start to listen to")

	serverCertificate = flag.String("cert", "", "public certificate shown by the server to it's clients")
	serverKey         = flag.String("key", "", "the server's certificate private key")
	caCertificate     = flag.String("ca-cert", "", "the CA certificate")
)

func main() {
	flag.Parse()

	var client *http.Client
	if *caCertificate == "" {
		client = http.DefaultClient
	} else {
		client = internal.HttpsClient(*caCertificate)
	}
	probeService := service.NewProbeService(client)
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
		Addr: "0.0.0.0:" + *port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if *serverCertificate == "" || *serverKey == "" {
			log.Printf("Starting HTTP server on port %s...\n", *port)
			if err := srv.ListenAndServe(); err != nil {
				log.Println(err)
			}
		} else {
			log.Printf("Starting HTTPs server on port %s...\n", *port)
			if err := srv.ListenAndServeTLS(*serverCertificate, *serverKey); err != nil {
				log.Println(err)
			}
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), *wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	_ = srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

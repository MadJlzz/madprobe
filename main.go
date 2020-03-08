package main

import (
	"context"
	"flag"
	"github.com/gorilla/mux"
	"github.com/madjlzz/madprobe/controller"
	"github.com/madjlzz/madprobe/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	var port string
	flag.StringVar(&port, "port", "3000", "the port for which the server will start to listen to")
	flag.Parse()

	probeService := service.NewProbeService()
	probeController := controller.NewProbeController(probeService)

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/probe/create", probeController.Create).
		Methods(http.MethodPost)

	srv := &http.Server{
		Addr: "0.0.0.0:" + port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Printf("Starting our server on port %s...\n", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
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

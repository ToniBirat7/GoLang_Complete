package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/birat/restapi/internal/config"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// Setup router
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Rest API"))
	})

	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("server started at", slog.String("address", cfg.Addr))
		err := server.ListenAndServe()

		if err != nil {
			log.Fatalf("failed to start server")
		}
	}()

	// Blocking, Until the Channel is not Notified by Signal
	<-done // We are taking out something from the channel

	// Main goroutine executes the code as usual
	slog.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx) // If the server cannot shutdown in the given time returns error

	if err != nil {
		slog.Error("failed to shutdown", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")
}

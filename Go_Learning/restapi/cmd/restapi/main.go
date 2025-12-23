package main

import (
	"fmt"
	"log"
	"net/http"

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

	fmt.Printf("Server Started %s\n", cfg.HTTPServer.Addr)
	err := server.ListenAndServe()

	if err != nil {
		log.Fatalf("failed to start server")
	}
}
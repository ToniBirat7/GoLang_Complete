package main

import (
	"net/http"

	"github.com/birat/restapi/internal/config"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// Setup router
	router := http.ServeMux()
}

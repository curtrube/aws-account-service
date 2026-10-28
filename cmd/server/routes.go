package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /v1/health", app.healthCheckHandler)

	// Accounts
	//mux.HandleFunc("GET /v1/account", app.getAccountHandler)

	return mux
}

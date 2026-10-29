package main

import (
	"log/slog"
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "ok",
		"environment": app.cfg.Environment,
		"version":     version,
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		slog.Error(err.Error())
		http.Error(
			w,
			"The server encountered a problem and could not process your request",
			http.StatusInternalServerError,
		)
	}
}

func (app *application) listAccountsHandler(w http.ResponseWriter, r *http.Request) {
	accounts, err := app.repo.List(r.Context())
	if err != nil {
		slog.Error("list accounts", "error", err)
		http.Error(
			w,
			"The server encoutered a problem and could not process your request",
			http.StatusInternalServerError,
		)
	}

	data := map[string]any{
		"accounts": accounts,
	}
	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		slog.Error("json encoding", "error", err)
		http.Error(
			w,
			"The server encountered a problem and could not process your request",
			http.StatusInternalServerError,
		)
	}
}

//func (app *application) getAccountHandler(w http.ResponseWriter, r *http.Request) { }

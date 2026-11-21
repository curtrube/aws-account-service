package main

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/curtrube/aws-account-service/internal/account"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"environment": app.cfg.Environment,
		"status":      "ok",
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
	accounts, err := app.account.List(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	data := envelope{"accounts": accounts}

	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getAccountHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 0 {
		http.NotFound(w, r)
		return
	}

	acct, err := app.account.Get(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, account.ErrNoAccount):
			http.NotFound(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	data := envelope{"accounts": acct}

	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

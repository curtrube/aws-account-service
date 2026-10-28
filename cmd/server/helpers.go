package main

import (
	"encoding/json"
	"net/http"
)

// writeJSON is a helper for sending JSON responses.
//
// Accepts a http.ResponseWriter, the http status code to send, the data to encode to JSON,
// and a headers map containing any additional HTTP headers to include in the response.
func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for key, values := range headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

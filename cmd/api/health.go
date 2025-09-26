package main

import "net/http"

func (app *Application) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("ping"))
}

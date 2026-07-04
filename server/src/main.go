package main

import (
	"effectivemobiletesttask/apiRegistration"
	"effectivemobiletesttask/logger"
	"net/http"
)

func main() {
	apiRegistration.Register()
	mainLog := logger.New("/health")
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/health",
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				mainLog.Print("Got request: ", r)
				if r.Method != "GET" {
					mainLog.Print("Got wrong request type")
				}
				w.WriteHeader(200)
				w.Write([]byte("Health is OK"))
			},
		),
	)
	mainLog.Print(
		http.ListenAndServe(
			"0.0.0.0:8080",
			mux,
		),
	)
}

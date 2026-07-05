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
					mainLog.Print("Got wrong request type:", r.Method, "instead of GET")
				}
				w.WriteHeader(200)
				_, err := w.Write([]byte("Health is OK"))
				if err != nil {
					mainLog.Print("During response got error:", err)
				}
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

package main

import (
	"effectivemobiletesttask/apiRegistration"
	"effectivemobiletesttask/db"
	"effectivemobiletesttask/server_logger"
	"net/http"
)

func main() {
	apiRegistration.Register()
	mainLog := server_logger.New("/health")
	db.Get()
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/health",
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				mainLog.Print("Got request: ", r)
				if r.Method != "GET" {
					mainLog.Print("Got wrong request type:", r.Method, "instead of GET")
				} else {
					w.WriteHeader(200)
					_, err := w.Write([]byte("Health is OK"))
					if err != nil {
						mainLog.Print("During response got error:", err)
					}
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

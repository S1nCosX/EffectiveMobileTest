package health_handlers

import (
	"effectivemobiletesttask/handlers"
	"effectivemobiletesttask/server_logger"
	"net/http"
)

func AddInMux(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		handlers.WriteResponse(w, server_logger.New("Health handler"), http.StatusOK, "Health is OK")
	})
}

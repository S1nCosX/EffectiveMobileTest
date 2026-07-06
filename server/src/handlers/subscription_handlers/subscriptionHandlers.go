package subscription_handlers

import (
	"effectivemobiletesttask/handlers"
	"effectivemobiletesttask/server_logger"
	"net/http"
	"strconv"
)

var logger = server_logger.New("Subscription handlers")

func handleCreate(w http.ResponseWriter, r *http.Request) {
	handlers.WriteResponse(w, logger, http.StatusOK, "got create req")
}

func handleRead(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, "Got read request with non-numerical id")
		return
	}
	if id < 0 {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, "Got read request with id less then 0")
		return
	}
	handlers.WriteResponse(w, logger, http.StatusOK, "Got correct read request")
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, "Got update request with non-numerical id")
		return
	}
	if id < 0 {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, "Got update request with id less then 0")
		return
	}

	handlers.WriteResponse(w, logger, http.StatusOK, "Got correct update request")
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, "Got delete request with non-numerical id")
		return
	}
	if id < 0 {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, "Got delete request with id less then 0")
		return
	}
	handlers.WriteResponse(w, logger, http.StatusOK, "Got correct delete request")
}

func handleList(w http.ResponseWriter, r *http.Request) {
	handlers.WriteResponse(w, logger, http.StatusOK, "Got list request")
}

func handleSummaryInPeriod(w http.ResponseWriter, r *http.Request) {
	handlers.WriteResponse(w, logger, http.StatusOK, "Got sum in period request")
}

func AddInMux(mux *http.ServeMux) {
	mux.HandleFunc("POST /subscriptions/create", handleCreate)
	mux.HandleFunc("DELETE /subscriptions/delete/{id}", handleDelete)
	mux.HandleFunc("GET /subscriptions/read/{id}", handleRead)
	mux.HandleFunc("PUT /subscriptions/update/{id}", handleUpdate)
	mux.HandleFunc("GET /subscriptions/list", handleList)

	mux.HandleFunc("GET /subscriptions/get_sum", handleSummaryInPeriod)
}

package subscription_handlers

import (
	"database/sql"
	subscriptions_service "effectivemobiletesttask/db/table_services/subscription_service"
	"effectivemobiletesttask/handlers"
	"effectivemobiletesttask/server_logger"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

var logger = server_logger.New("Subscription handlers")

func handleCreate(w http.ResponseWriter, r *http.Request) {
	var subscription subscriptions_service.Subscription

	err := json.NewDecoder(r.Body).Decode(&subscription)
	if err != nil {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, fmt.Sprintf("BODY DESERIALIZATION ERROR : %s", err))
		return
	}

	id, err := subscriptions_service.Create(subscription)
	if err != nil {
		handlers.WriteResponse(w, logger, http.StatusInternalServerError, fmt.Sprintf("CREATE ERROR: %s", err))
		return
	}
	handlers.WriteResponse(w, logger, http.StatusOK, fmt.Sprintf("Added new subscription. ID: %d", id))
}

func handleRead(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, "Non-numerical ID")
		return
	}
	if id < 0 {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, "ID less then 0")
		return
	}

	subscription, err := subscriptions_service.Read(uint(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			handlers.WriteResponse(w, logger, http.StatusNotFound, "No rows found")
		} else {
			handlers.WriteResponse(w, logger, http.StatusInternalServerError, fmt.Sprintf("READING ERROR: %s", err))
		}
		return
	}

	str, err := json.Marshal(subscription)
	if err != nil {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, fmt.Sprintf("SERIALIZATION ERROR: %s", err))
		return
	}
	handlers.WriteResponse(w, logger, http.StatusOK, string(str))
}

func parseSubscriptionMap(subscription *subscriptions_service.Subscription) (ret map[string]string) {
	ret = make(map[string]string)

	if subscription.ServiceName != "" {
		ret["service_name"] = subscription.ServiceName
	}

	if subscription.UserId != "" {
		ret["user_id"] = subscription.UserId
	}

	if subscription.Price >= 0 {
		ret["price"] = strconv.Itoa(int(subscription.Price))
	}

	if subscription.StartDate != "" {
		ret["start_date"] = subscription.StartDate
	}

	if subscription.EndDate != nil {
		ret["end_date"] = subscription.StartDate
	}

	return ret
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, "Non-numerical id")
		return
	}
	if id < 0 {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, "ID less then 0")
		return
	}

	var subscription subscriptions_service.Subscription

	err = json.NewDecoder(r.Body).Decode(&subscription)
	if err != nil {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, fmt.Sprintf("BODY DESERIALIZATION ERROR : %s", err))
		return
	}

	replaces := parseSubscriptionMap(&subscription)

	new_subscription, err := subscriptions_service.Update(uint(id), replaces)

	if err != nil {
		handlers.WriteResponse(w, logger, http.StatusInternalServerError, fmt.Sprintf("UPDATE ERROR: %s", err))
		return
	}

	str, err := json.Marshal(new_subscription)
	if err != nil {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, fmt.Sprintf("SERIALIZATION ERROR: %s", err))
		return
	}
	handlers.WriteResponse(w, logger, http.StatusOK, fmt.Sprintf("Updated: %s", str))
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, "Non-numerical id")
		return
	}
	if id < 0 {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, "ID less then 0")
		return
	}

	err = subscriptions_service.Delete(uint(id))
	if err != nil {
		handlers.WriteResponse(w, logger, http.StatusInternalServerError, fmt.Sprintf("DELETE error: %s", err))
	}
	handlers.WriteResponse(w, logger, http.StatusOK, "Deleted")
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

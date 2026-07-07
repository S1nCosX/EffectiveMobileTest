package subscription_handlers

import (
	"database/sql"
	"effectivemobiletesttask/db/dto"
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

func parseSubscriptionMap(subscription *dto.SubscriptionReadDTO) (ret map[string]string) {
	ret = make(map[string]string)

	if subscription.ServiceName != nil {
		ret["service_name"] = *subscription.ServiceName
	}

	if subscription.UserId != nil {
		ret["user_id"] = *subscription.UserId
	}

	if subscription.Price != nil {
		ret["price"] = strconv.Itoa(int(*subscription.Price))
	}

	if subscription.StartDate != nil {
		ret["start_date"] = "01-" + *subscription.StartDate
	}

	if subscription.EndDate != nil {
		ret["end_date"] = "01-" + *subscription.EndDate
	}

	return ret
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	var subscription dto.SubscriptionReadDTO

	err := json.NewDecoder(r.Body).Decode(&subscription)
	if err != nil {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, fmt.Sprintf("BODY DESERIALIZATION ERROR : %s", err))
		return
	}

	inp := parseSubscriptionMap(&subscription)

	id, err := subscriptions_service.Create(&inp)
	if err != nil {
		handlers.WriteResponse(w, logger, http.StatusInternalServerError, fmt.Sprintf("DB ERROR: %s", err))
		return
	}
	handlers.WriteResponse(w, logger, http.StatusOK, strconv.FormatUint(uint64(id), 10))
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
			handlers.WriteResponse(w, logger, http.StatusInternalServerError, fmt.Sprintf("DB ERROR: %s", err))
		}
		return
	}

	resp, err := json.Marshal(subscription)
	if err != nil {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, fmt.Sprintf("SERIALIZATION ERROR: %s", err))
		return
	}
	handlers.WriteResponse(w, logger, http.StatusOK, string(resp))
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

	var subscription dto.SubscriptionReadDTO

	err = json.NewDecoder(r.Body).Decode(&subscription)
	if err != nil {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, fmt.Sprintf("DESERIALIZATION ERROR : %s", err))
		return
	}

	replaces := parseSubscriptionMap(&subscription)

	new_subscription, err := subscriptions_service.Update(uint(id), &replaces)

	if err != nil {
		handlers.WriteResponse(w, logger, http.StatusInternalServerError, fmt.Sprintf("DB ERROR: %s", err))
		return
	}

	resp, err := json.Marshal(new_subscription)
	if err != nil {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, fmt.Sprintf("SERIALIZATION ERROR: %s", err))
		return
	}
	handlers.WriteResponse(w, logger, http.StatusOK, fmt.Sprintf("Updated: %s", resp))
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
	if !r.URL.Query().Has("page") {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, "QUERY DOESNT HAVE page PARAMETER")
		return
	}

	page, err := strconv.ParseUint(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, fmt.Sprintf("CANNOT PARSE PAGE VALUE: %s", err))
		return
	}

	if !r.URL.Query().Has("page_size") {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, "QUERY DOESNT HAVE page_size PARAMETER")
		return
	}

	page_size, err := strconv.ParseUint(r.URL.Query().Get("page_size"), 10, 64)
	if err != nil {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, fmt.Sprintf("CANNOT PARSE PAGE SIZE VALUE: %s", err))
		return
	}

	ret, err := subscriptions_service.List(page, page_size)
	if err != nil {
		handlers.WriteResponse(w, logger, http.StatusInternalServerError, fmt.Sprintf("DB ERROR: %s", err))
		return
	}

	resp, err := json.Marshal(ret)
	if err != nil {
		handlers.WriteResponse(w, logger, http.StatusInternalServerError, fmt.Sprintf("SERIALIZATION ERROR: %s", err))
		return
	}

	handlers.WriteResponse(w, logger, http.StatusOK, string(resp))
}

func handleSummaryInPeriod(w http.ResponseWriter, r *http.Request) {

	var args []any
	if !r.URL.Query().Has("id") {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, "QUERY DOESN'T HAVE id PARAMETER")
		return
	}

	args = append(args, r.URL.Query().Get("id"))

	if !r.URL.Query().Has("service_name") {
		handlers.WriteResponse(w, logger, http.StatusBadRequest, "QUERY DOESN'T HAVE service_name PARAMETER")
		return
	}

	args = append(args, r.URL.Query().Get("service_name"))

	filter := "user_id = $1 AND service_name = $2"

	if r.URL.Query().Has("start_time") {
		filter += " AND start_date >= $3"

		args = append(args, "01-"+r.URL.Query().Get("start_time"))
		if r.URL.Query().Has("end_time") {
			filter += " AND end_date <= $4"

			args = append(args, "01-"+r.URL.Query().Get("end_time"))
		}
	} else {
		if r.URL.Query().Has("end_time") {
			filter += " AND end_date <= $3"

			args = append(args, "01-"+r.URL.Query().Get("end_time"))
		}
	}

	logger.Print("got sum with args: ", args)

	resp, err := subscriptions_service.GetSumInPeriod(&args, &filter)

	if err != nil {
		handlers.WriteResponse(w, logger, http.StatusInternalServerError, fmt.Sprintf("DB ERROR: %s", err))
		return
	}
	handlers.WriteResponse(w, logger, http.StatusOK, strconv.FormatUint(uint64(resp), 10))
}

func AddInMux(mux *http.ServeMux) {
	mux.HandleFunc("POST /subscriptions/create", handleCreate)
	mux.HandleFunc("DELETE /subscriptions/delete/{id}", handleDelete)
	mux.HandleFunc("GET /subscriptions/read/{id}", handleRead)
	mux.HandleFunc("PUT /subscriptions/update/{id}", handleUpdate)
	mux.HandleFunc("GET /subscriptions/list", handleList)

	mux.HandleFunc("GET /subscriptions/price_sum", handleSummaryInPeriod)
}

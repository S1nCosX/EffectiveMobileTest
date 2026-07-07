package main

import (
	"effectivemobiletesttask/apiRegistration"
	"effectivemobiletesttask/config"
	subscriptions_service "effectivemobiletesttask/db/table_services/subscription_service"
	_ "effectivemobiletesttask/docs"
	"effectivemobiletesttask/handlers/health_handlers"
	"effectivemobiletesttask/handlers/subscription_handlers"
	"effectivemobiletesttask/server_logger"
	"fmt"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Subscriptions API
// @host localhost:1337
// @BasePath /
func main() {
	apiRegistration.Register()
	logger := server_logger.New("main")

	mux := http.NewServeMux()

	health_handlers.AddInMux(mux)

	subscription_handlers.AddInMux(mux)
	subscriptions_service.Init()

	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)

	conf, err := config.Get()
	if err != nil {
		logger.Panic("Got errors during config reading", err)
	}

	logger.Print(http.ListenAndServe(fmt.Sprintf("%s:%d", conf.HOST, conf.PORT), mux))
}

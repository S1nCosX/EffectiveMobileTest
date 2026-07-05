package apiRegistration

import (
	"effectivemobiletesttask/config"
	"effectivemobiletesttask/logger"
	"fmt"
	"strconv"

	"github.com/hashicorp/consul/api"
)

func Register() {
	registerLogger := logger.New("Service Registration")
	conf, err := config.Get()
	if err != nil {
		registerLogger.Panic(err)
	}

	apiConfig := api.DefaultConfig()
	apiConfig.Address = getConsulAddr()

	client, err := api.NewClient(apiConfig)
	if err != nil {
		registerLogger.Panic("Problem with service registration in consul:", err)
	}
	registerLogger.Print("Consul client successfuly created")

	err = client.Agent().ServiceRegister(
		&api.AgentServiceRegistration{
			Name:    "server",
			ID:      conf.SERVICE_ID,
			Address: conf.SERVICE_ID,
			Port:    int(conf.PORT),
			Check: &api.AgentServiceCheck{
				HTTP:                           fmt.Sprintf("http://%s:%d/health", conf.SERVICE_ID, conf.PORT),
				Method:                         "GET",
				Timeout:                        "5s",
				Interval:                       "10s",
				DeregisterCriticalServiceAfter: "30s",
			},
		},
	)
	if err != nil {
		registerLogger.Panic("Got error during registration: ", err)
	}
	registerLogger.Print("Consul Registration succseed")
}

func getConsulAddr() string {
	conf, err := config.Get()
	if err != nil {
		logger.New("Consul address").Panic(err)
	}
	return conf.CONSUL_ADDR + ":" + strconv.FormatUint(uint64(conf.CONSUL_PORT), 10)
}

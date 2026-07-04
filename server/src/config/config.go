package config

import (
	"effectivemobiletesttask/logger"
	"os"
	"strconv"
	"sync"
)

type Config struct {
	CONSUL_ADDR string
	CONSUL_PORT uint16
	SERVICE_ID  string
	HOST        string
	PORT        uint16
}

var (
	instance Config
	once     sync.Once
	initErr  error
)

func Get() (*Config, error) {
	once.Do(
		func() {
			configLogger := logger.New("Config reader")
			serviceId, err := os.Hostname()
			if err != nil {
				configLogger.Panic("Failed to get service id as host name")
			}
			consulAddr, isExist := os.LookupEnv("CONSUL_ADDR")
			if !isExist {
				configLogger.Print("WARN: CONSUL_ADDR environment variable is not exist, using default value")
				consulAddr = "consul"
			}
			consulPort, isExist := os.LookupEnv("CONSUL_PORT")
			if !isExist {
				configLogger.Print("WARN: CONSUL_PORT environment variable is not exist, using default value")
				consulPort = "8500"
			}
			consulPortInt, err := strconv.ParseUint(consulPort, 10, 16)
			if err != nil {
				configLogger.Panic("During port parsing got error:\"", err, "\"")
			}

			host, isExist := os.LookupEnv("HOST")
			if !isExist {
				configLogger.Print("WARN: HOST environment variable is not exist, using default value")
				host = "0.0.0.0"
			}
			port, isExist := os.LookupEnv("PORT")
			if !isExist {
				configLogger.Print("WARN: PORT environment variable is not exist, using default value")
				port = "8080"
			}
			portInt, err := strconv.ParseUint(port, 10, 16)
			if err != nil {
				configLogger.Panic("During port parsing got error:\"", err, "\"")
			}

			instance = Config{
				SERVICE_ID:  serviceId,
				CONSUL_ADDR: consulAddr,
				CONSUL_PORT: uint16(consulPortInt),
				HOST:        host,
				PORT:        uint16(portInt),
			}
		},
	)
	return &instance, initErr
}

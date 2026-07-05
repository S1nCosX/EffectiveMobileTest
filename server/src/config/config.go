package config

import (
	"effectivemobiletesttask/server_logger"
	"log"
	"os"
	"strconv"
	"sync"
)

type Config struct {
	// Consul info
	CONSUL_ADDR string
	CONSUL_PORT uint16
	// Server instance id
	SERVICE_ID string
	// Server addr
	HOST string
	PORT uint16
	// Psql data
	DB_HOST string
	DB_PORT uint16
	DB_NAME string
	DB_USER string
	DB_PW   string
}

var (
	instance Config
	once     sync.Once
	initErr  error
	logger   *log.Logger
)

func getStrEnv(key string, defaultValue string) string {
	value, isExist := os.LookupEnv(key)
	if !isExist {
		logger.Printf("WARN: %s environment variable is not exist, using default value", key)
		return defaultValue
	}
	return value
}

func getUint16(key string, defaultValue uint16) uint16 {
	strValue, isExist := os.LookupEnv(key)
	if !isExist {
		logger.Printf("WARN: %s environment variable is not exist, using default value", key)
		return defaultValue
	}

	value, err := strconv.ParseUint(strValue, 10, 16)
	if err != nil {
		logger.Panicf("During port parsing got error: %s ", err)
	}
	return uint16(value)
}

func (*Config) init() {
	logger = server_logger.New("Config reader")

	serviceId, err := os.Hostname()
	if err != nil {
		logger.Panic("Failed to get service id as host name")
	}

	bytes, err := os.ReadFile(getStrEnv("DB_PW", ""))

	if err != nil {
		logger.Printf("Tryed to get secret, but got error: %s", err)
		bytes = []byte("postgres")
	}

	pw := string(bytes)

	instance = Config{
		SERVICE_ID:  serviceId,
		CONSUL_ADDR: getStrEnv("CONSUL_ADDR", "consul"),
		CONSUL_PORT: getUint16("CONSUL_PORT", 8500),
		HOST:        getStrEnv("HOST", "0.0.0.0"),
		PORT:        getUint16("PORT", 8080),
		DB_HOST:     getStrEnv("DB_HOST", "db"),
		DB_PORT:     getUint16("DB_PORT", 5432),
		DB_NAME:     getStrEnv("DB_NAME", "postgres"),
		DB_USER:     getStrEnv("DB_USER", "postgres"),
		DB_PW:       pw,
	}

	logger.Print(instance)
}

func Get() (*Config, error) {
	once.Do(instance.init)
	return &instance, initErr
}

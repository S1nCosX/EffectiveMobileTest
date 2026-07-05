package db

import (
	"database/sql"
	"effectivemobiletesttask/config"
	"effectivemobiletesttask/server_logger"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

type DatabaseDriver struct {
	conn_str string
	conn     *sql.DB
	once     sync.Once
}

var (
	once     sync.Once
	instance DatabaseDriver
	initErr  error
	logger   *log.Logger
)

func (*DatabaseDriver) init() {
	logger = server_logger.New("Database driver")

	conf, err := config.Get()
	if err != nil {
		logger.Panic("Config initiated  with error")
	}

	instance.conn_str = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.DB_HOST, conf.DB_PORT, conf.DB_USER, conf.DB_PW, conf.DB_NAME)

	instance.conn, err = sql.Open("postgres", instance.conn_str)
	if err != nil {
		logger.Panicf("Got error during db connection: %s", err)
	}
}

func Get() (*DatabaseDriver, error) {
	once.Do(instance.init)
	return &instance, initErr
}

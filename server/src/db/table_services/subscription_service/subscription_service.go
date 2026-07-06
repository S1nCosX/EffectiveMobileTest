package subscriptions_service

import (
	db "effectivemobiletesttask/db/driver"
	"effectivemobiletesttask/server_logger"
	"log"
)

type Subscription struct {
	Id          uint
	ServiceName string  `json:"service_name"`
	UserId      string  `json:"user_id"`
	Price       uint    `json:"price"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

var (
	driver *db.DatabaseDriver
	err    error
	logger *log.Logger
)

func Init() {
	driver, err = db.Get()
	logger = server_logger.New("Subscription service")

	if err != nil {
		logger.Panic("Init error. Database driver getting got err: ", err)
	}
}

func Create(newSub Subscription) (uint, error) {
	var id uint
	if newSub.EndDate != nil {
		err := driver.Conn.QueryRow("INSERT INTO subscriptions (service_name, user_id, price, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			newSub.ServiceName,
			newSub.UserId,
			newSub.Price,
			newSub.StartDate,
			newSub.EndDate,
		).Scan(&id)
		return id, err
	} else {
		err := driver.Conn.QueryRow("INSERT INTO subscriptions (service_name, user_id, price, start_date) VALUES ($1, $2, $3, $4) RETURNING id",
			newSub.ServiceName,
			newSub.UserId,
			newSub.Price,
			newSub.StartDate,
		).Scan(&id)
		return id, err
	}
}

func Read(id uint) (Subscription, error) {
	var ret Subscription
	err := driver.Conn.QueryRow("SELECT id, service_name, user_id, price, to_char(start_date, 'MM-YYYY') as start_date, to_char(end_date, 'MM-YYYY') as end_date FROM subscriptions WHERE id = $1", id).Scan(
		&ret.Id,
		&ret.ServiceName,
		&ret.UserId,
		&ret.Price,
		&ret.StartDate,
		&ret.EndDate,
	)
	return ret, err
}

func Update(id uint, replaces map[string]string) (Subscription, error) {
	var fields []rune
	var args []rune
	var args_values []string

	i := 0
	for k, v := range replaces {
		fields = append(fields, []rune(k+",")...)
		args = append(fields, []rune(string(i)+",")...)
		args_values = append(args_values)
		i++
	}

	var ret Subscription
	err := driver.Conn.QueryRow("UPDATE subscriptions SET (%s) = (%s) WHERE id = $%d RETURNING NEW", id).Scan(
		&ret.Id,
		&ret.ServiceName,
		&ret.UserId,
		&ret.Price,
		&ret.StartDate,
		&ret.EndDate,
	)
	return ret, err
}

func Delete(id uint) error {
	_, err := driver.Conn.Exec("DELETE FROM subscriptions WHERE id = $1", id)
	return err
}

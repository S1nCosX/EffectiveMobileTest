package subscriptions_service

import (
	db "effectivemobiletesttask/db/driver"
	"effectivemobiletesttask/db/dto"
	"effectivemobiletesttask/server_logger"
	"fmt"
	"log"
	"strconv"
	"strings"
)

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

func convertMapToQueryObjects(mp *map[string]string) (variables string, args string, values []any) {
	var variables_array []string
	var args_array []string

	i := 1
	for k, v := range *mp {
		variables_array = append(variables_array, k)
		args_array = append(args_array, "$"+strconv.Itoa(i))
		values = append(values, v)
		i++
	}
	variables = strings.Join(variables_array, ", ")
	args = strings.Join(args_array, ", ")
	return variables, args, values
}

func Create(newSub *map[string]string) (uint, error) {
	var id uint

	var fields, args, args_values = convertMapToQueryObjects(newSub)
	query := fmt.Sprintf("INSERT INTO subscriptions (%s) VALUES (%s) RETURNING id", fields, args)
	err := driver.Conn.QueryRow(query, args_values...).Scan(&id)
	return id, err

}

func Read(id uint) (dto.SubscriptionDTO, error) {
	query :=
		`SELECT
		id,
		service_name,
		user_id,
		price,
		to_char(start_date, 'MM-YYYY') as start_date,
		to_char(end_date, 'MM-YYYY') as end_date
	FROM subscriptions WHERE id = $1`
	var ret dto.SubscriptionDTO
	err := driver.Conn.QueryRow(query, id).Scan(
		&ret.Id,
		&ret.ServiceName,
		&ret.UserId,
		&ret.Price,
		&ret.StartDate,
		&ret.EndDate,
	)
	return ret, err
}

func Update(id uint, replaces *map[string]string) (dto.SubscriptionDTO, error) {
	var fields, args, args_values = convertMapToQueryObjects(replaces)

	query := fmt.Sprintf(
		`UPDATE subscriptions
		SET (%s) = (%s)
		WHERE id = $%d
		RETURNING
			new.id,
			new.service_name,
			new.user_id,
			new.price,
			to_char(new.start_date, 'MM-YYYY') as start_date,
			to_char(new.end_date, 'MM-YYYY') as end_date`,
		fields, args, len(args_values)+1)

	args_values = append(args_values, id)

	var ret dto.SubscriptionDTO
	err := driver.Conn.QueryRow(query, args_values...).Scan(
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

func List(page uint64, page_size uint64) (ret []dto.SubscriptionDTO, err error) {
	query :=
		`SELECT id,
		service_name,
		user_id,
		price,
		to_char(start_date, 'MM-YYYY') as start_date,
		to_char(end_date, 'MM-YYYY') as end_date 
	FROM subscriptions
	LIMIT $1
	OFFSET $2 ROWS`
	rows, err := driver.Conn.Query(query, page_size, page*page_size)

	for rows.Next() {
		var row dto.SubscriptionDTO
		rows.Scan(
			&row.Id,
			&row.ServiceName,
			&row.UserId,
			&row.Price,
			&row.StartDate,
			&row.EndDate,
		)
		ret = append(ret, row)
	}
	return ret, err
}

func GetSumInPeriod(args *[]any, filter *string) (ret uint, err error) {
	query := fmt.Sprintf("SELECT sum(price) FROM subscriptions WHERE %s", *filter)

	err = driver.Conn.QueryRow(query, *args...).Scan(&ret)

	return ret, err
}

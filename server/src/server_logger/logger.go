package server_logger

import "log"

func New(prefix string) (logger *log.Logger) {
	logger = log.New(
		log.Writer(),
		"["+prefix+"]: ",
		log.Ldate+log.LUTC,
	)
	return logger
}

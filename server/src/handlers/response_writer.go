package handlers

import (
	"log"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, logger *log.Logger, code int, message string) {
	logger.Print(message)
	w.WriteHeader(code)
	w.Write([]byte(message))
}

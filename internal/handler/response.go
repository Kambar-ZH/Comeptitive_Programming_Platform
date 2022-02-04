package handler

import (
	"encoding/json"
	"net/http"
	"site/internal/logger"
)

func Error(w http.ResponseWriter, r *http.Request, code int, err error) {
	logger.Logger.Error(err.Error())
	Respond(w, r, code, map[string]string{"error": err.Error()})
}

func Respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

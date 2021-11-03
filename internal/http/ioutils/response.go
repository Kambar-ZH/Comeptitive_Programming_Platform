package ioutils

import (
	"encoding/json"
	"log"
	"net/http"
)

func Error(w http.ResponseWriter, r *http.Request, code int, err error) {
	log.Println(err)
	Respond(w, r, code, map[string]string{"error": err.Error()})
}

func Respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
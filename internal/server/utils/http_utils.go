package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteString(w http.ResponseWriter, message string) {
	_, err := w.Write([]byte(message))
	if err != nil {
		log.Fatal("Error while writing response.")
	}
}

func WriteMessageResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.WriteHeader(httpStatusCode)
	WriteString(w, message)
}

func WriteJsonResponse(w http.ResponseWriter, object interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(object)
	if err != nil {
		log.Fatal("Error while writing json response.")
	}
}

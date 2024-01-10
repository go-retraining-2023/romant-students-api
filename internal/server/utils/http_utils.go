package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
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

func RetrieveStudentId(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	// getting student id from url
	idString := r.URL.Path[len("/api/v1/students/"):]
	studentId, err := uuid.Parse(idString)
	if err != nil {
		http.Error(w, "Ivalid student identifier", http.StatusBadRequest)
		return uuid.UUID{}, err
	}

	return studentId, nil
}

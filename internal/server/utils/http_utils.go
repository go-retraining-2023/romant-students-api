package server

import (
	"net/http"

	"github.com/google/uuid"
)

func WriteString(w http.ResponseWriter, message string) {
	w.Write([]byte(message))
}

func WriteError(w http.ResponseWriter, message string, httpStatusCode int) {
	w.WriteHeader(httpStatusCode)
	WriteString(w, message)
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

package server

import (
	"net/http"

	"github.com/google/uuid"
)

func WriteString(w http.ResponseWriter, value string) {
	w.Write([]byte(value))
}

func RetrieveStudentId(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	// getting student id from url
	idString := r.URL.Path[len("/api/students/"):]
	studentId, err := uuid.Parse(idString)
	if err != nil {
		http.Error(w, "Ivalid student identifier", http.StatusBadRequest)
		return uuid.UUID{}, err
	}

	return studentId, nil
}

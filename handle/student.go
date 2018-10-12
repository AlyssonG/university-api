package handle

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/alyssong/university-api/storage"
)

//StudentHandle defines the API endpoints implementations for student operations
type StudentHandle struct {
	Logger  *log.Logger
	Storage storage.Student
}

//GetStudent expects a parameter ID and returns data for a student with that id
func (sh *StudentHandle) GetStudent(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	value := query.Get("id")
	studentID, err := strconv.Atoi(value)
	if err != nil {
		http.Error(w, "invalid type for student id. expecting int", http.StatusInternalServerError)
		sh.Logger.Println("invalid type for student id", "value", value)
		return
	}

	student, err := sh.Storage.Get(studentID)
	if err != nil || student == nil {
		http.Error(w, "student nostudent not foundt found", http.StatusNotFound)
		sh.Logger.Println("student not found for id", studentID, "err", err)
		return
	}

	body, err := json.Marshal(student)
	if err != nil {
		http.Error(w, "error while enconding student", http.StatusInternalServerError)
		sh.Logger.Println("error while encoding student to json", "err", err)
	}
	w.Header().Set("Content-Type", "Application/json")
	w.Write(body)
}

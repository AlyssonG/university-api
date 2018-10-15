package handle

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/alyssong/university-api/university"

	"github.com/alyssong/university-api/storage"
)

//StudentHandle defines the API endpoints implementations for student operations
type StudentHandle struct {
	Logger  *log.Logger
	Storage storage.Student
}

//GetStudent expects a parameter ID and returns data for a student with that id
func (sh *StudentHandle) GetStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	studentID := query.Get("id")

	if studentID != "" {
		sh.getSpecificStudent(studentID, w, r)
	} else {
		sh.getStudents(w, r)
	}
}

func (sh *StudentHandle) getSpecificStudent(studentID string, w http.ResponseWriter, r *http.Request) {
	student, err := sh.Storage.Get(studentID)
	if err != nil || student == nil {
		http.Error(w, "student not foundt found", http.StatusNotFound)
		sh.Logger.Println("student not found for id", studentID, "err", err)
		return
	}

	body, err := json.Marshal(student)
	if err != nil {
		http.Error(w, "error while enconding student", http.StatusInternalServerError)
		sh.Logger.Println("error while encoding student to json", "err", err)
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	w.Write(body)
}

func (sh *StudentHandle) getStudents(w http.ResponseWriter, r *http.Request) {
	students, err := sh.Storage.GetAll()
	if err != nil {
		http.Error(w, "error while recovering students", http.StatusInternalServerError)
		sh.Logger.Println("error while recovering students", "err", err)
		return
	}

	body, err := json.Marshal(students)
	if err != nil {
		http.Error(w, "error while enconding students", http.StatusInternalServerError)
		sh.Logger.Println("error while encoding students to json", "err", err)
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	w.Write(body)
}

func (sh *StudentHandle) SetStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, "invalid request body", http.StatusInternalServerError)
		sh.Logger.Println("invalid body to create student", "body", string(body))
		return
	}

	student := &university.Student{}
	err = json.Unmarshal(body, student)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusInternalServerError)
		sh.Logger.Println("error while creating student object", "body", string(body))
		return
	}

	if s, _ := sh.Storage.Get(student.Code); s != nil {
		http.Error(w, "a student with this data is registred in the system", http.StatusConflict)
		sh.Logger.Println("a student with this data is registred in the system", "code", s.Code)
		return
	}

	_, err = sh.Storage.Set(student)
	if err != nil {
		http.Error(w, "error while saving student in db", http.StatusInternalServerError)
		sh.Logger.Println("error while saving student in db", "err", err)
		return
	}

	w.Write([]byte("created"))
}

func (sh *StudentHandle) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	code := query.Get("code")

	student, _ := sh.Storage.Get(code)
	if student == nil {
		http.Error(w, "this student does not exists", http.StatusNotFound)
		return
	}

	err := sh.Storage.Delete(code)
	if err != nil {
		http.Error(w, "error while deleting record", http.StatusInternalServerError)
		sh.Logger.Println("error while deleting student record", "err", err)
		return
	}

	w.Write([]byte("deleted"))
}

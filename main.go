package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alyssong/university-api/storage/memory"
	"github.com/alyssong/university-api/university"

	"github.com/alyssong/university-api/handle"
)

func main() {
	logger := log.New(os.Stdout, "-- University API --", 0)
	Student := handle.Student{
		Logger: logger,
		Storage: &memory.StudentStorage{
			Store: make(map[string]*university.Student),
		},
	}

	professorHandle := handle.Professor{
		Logger: logger,
		Storage: &memory.ProfessorStorage{
			Store: make(map[string]*university.Professor),
		},
	}

	http.HandleFunc("/student", Student.GetStudent)
	http.HandleFunc("/student/new", Student.SetStudent)
	http.HandleFunc("/student/delete", Student.DeleteStudent)
	http.HandleFunc("/student/update", Student.UpdateStudent)

	http.HandleFunc("/professor", professorHandle.GetProfessor)
	http.HandleFunc("/professor/new", professorHandle.SetProfessor)
	http.HandleFunc("/professor/delete", professorHandle.DeleteProfessor)
	http.HandleFunc("/professor/update", professorHandle.UpdateProfessor)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server error", "err", err)
	}
}

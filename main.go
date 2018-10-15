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
	studentHandle := handle.StudentHandle{
		Logger: log.New(os.Stdout, "-- University API --", 0),
		Storage: &memory.StudentStorage{
			Store: make(map[string]*university.Student),
		},
	}

	http.HandleFunc("/student", studentHandle.GetStudent)
	http.HandleFunc("/student/new", studentHandle.SetStudent)
	http.HandleFunc("/student/delete", studentHandle.DeleteStudent)
	http.HandleFunc("/student/update", studentHandle.UpdateStudent)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server error", "err", err)
	}
}

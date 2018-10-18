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
	logger := log.New(os.Stdout, "-- University API -- ", log.LstdFlags)
	student := handle.Student{
		Logger: logger,
		Storage: &memory.StudentStorage{
			Store: make(map[string]*university.Student),
		},
	}

	professor := handle.Professor{
		Logger: logger,
		Storage: &memory.ProfessorStorage{
			Store: make(map[string]*university.Professor),
		},
	}

	http.HandleFunc("/student", student.Handle)
	http.HandleFunc("/professor", professor.Handle)

	logger.Println("Server started with success")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server error", "err", err)
	}
}

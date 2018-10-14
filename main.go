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

	student := &university.Student{
		Name: "Rogerinho",
		Code: "ABC",
	}

	studentHandle.Storage.Set(student)

	http.HandleFunc("/student", studentHandle.GetStudent)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server error", "err", err)
	}
}

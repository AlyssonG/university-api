package memory_test

import (
	"testing"

	"github.com/alyssong/university-api/university"

	"github.com/alyssong/university-api/storage/memory"
)

func TestStorageSet(t *testing.T) {
	storage := memory.StudentStorage{
		Store: make(map[string]*university.Student),
	}

	student := &university.Student{
		Name: "Rogerinho",
		Code: "ABC",
	}

	code, err := storage.Set(student)
	if err != nil {
		t.Fatal("this operation does not expect error")
	}

	if code != student.Code {
		t.Error("unexpected Code", "code", code, "expected", student.Code)
	}

	_, err = storage.Set(nil)
	if err == nil {
		t.Error("set method cannot accept nil values without returning an error")
	}
}

func TestStorageGet(t *testing.T) {
	storage := memory.StudentStorage{
		Store: make(map[string]*university.Student),
	}

	student := &university.Student{
		Name: "Rogerinho",
		Code: "ABC",
	}

	storage.Store["ABC"] = student

	storedStudent, err := storage.Get("ABC")
	if err != nil {
		t.Fatal("unexpected error", "err", err)
	}

	if storedStudent == nil {
		t.Fatal("invalid return for storage get. student is nil")
	}

	if student.Name != storedStudent.Name || student.Code != storedStudent.Code {
		t.Error("invalid data for stored student")
	}
}

func TestStorageDelete(t *testing.T) {
	storage := memory.StudentStorage{
		Store: make(map[string]*university.Student),
	}

	student := &university.Student{
		Name: "Rogerinho",
		Code: "Do Ingá",
	}

	storage.Store["ABC"] = student

	err := storage.Delete("ABC")
	if err != nil {
		t.Error("delete operation is not ok", "err", err)
	}

	if storage.Store["ABC"] != nil {
		t.Error("delete operation is not deleting from memory")
	}
}

func TestGetAll(t *testing.T) {
	storage := memory.StudentStorage{
		Store: make(map[string]*university.Student),
	}

	studentA := &university.Student{
		Name: "Rogerinho",
		Code: "ABC",
	}

	storage.Store["ABC"] = studentA

	studentB := &university.Student{
		Name: "Rogerinho",
		Code: "CBA",
	}

	storage.Store["CBA"] = studentB

	students := []*university.Student{
		studentA,
		studentB,
	}

	result, _ := storage.GetAll()
	for _, expected := range result {
		found := false
		for _, student := range students {
			if student.Code == expected.Code {
				found = true
			}
		}

		if !found {
			t.Error("it is missing a student in get all return")
		}
		return
	}
}

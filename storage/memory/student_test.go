package memory_test

import (
	"testing"

	"github.com/alyssong/university-api/university"

	"github.com/alyssong/university-api/storage/memory"
)

func TestStorageSet(t *testing.T) {
	storage := memory.StudentStorage{
		Store: make(map[int]*university.Student),
	}

	student := &university.Student{
		ID:   1,
		Name: "Rogerinho",
		Code: "Do Ingá",
	}

	id, err := storage.Set(student)
	if err != nil {
		t.Fatal("this operation does not expect error")
	}

	if id != student.ID {
		t.Error("unexpected ID", "id", id, "expected", student.ID)
	}

	_, err = storage.Set(nil)
	if err == nil {
		t.Error("set method cannot accept nil values without returning an error")
	}
}

func TestStorageGet(t *testing.T) {
	storage := memory.StudentStorage{
		Store: make(map[int]*university.Student),
	}

	student := &university.Student{
		ID:   1,
		Name: "Rogerinho",
		Code: "Do Ingá",
	}

	storage.Store[1] = student

	storedStudent, err := storage.Get(1)
	if err != nil {
		t.Fatal("unexpected error", "err", err)
	}

	if storedStudent == nil {
		t.Fatal("invalid return for storage get. student is nil")
	}

	if student.ID != storedStudent.ID || student.Name != storedStudent.Name || student.Code != storedStudent.Code {
		t.Error("invalid data for stored student")
	}
}

func TestStorageDelete(t *testing.T) {
	storage := memory.StudentStorage{
		Store: make(map[int]*university.Student),
	}

	student := &university.Student{
		ID:   1,
		Name: "Rogerinho",
		Code: "Do Ingá",
	}

	storage.Store[1] = student

	err := storage.Delete(1)
	if err != nil {
		t.Error("delete operation is not ok", "err", err)
	}

	if storage.Store[1] != nil {
		t.Error("delete operation is not deleting from memory")
	}
}

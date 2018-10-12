package memory_test

import (
	"testing"

	"github.com/alyssong/university-api/university"

	"github.com/alyssong/university-api/storage/memory"
)

func TestStorageSet(t *testing.T) {
	storage := memory.StudentStorage{}
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
}

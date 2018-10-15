package memory_test

import (
	"testing"

	"github.com/alyssong/university-api/university"

	"github.com/alyssong/university-api/storage/memory"
)

func TestProfessorStorageSet(t *testing.T) {
	storage := memory.ProfessorStorage{
		Store: make(map[string]*university.Professor),
	}

	professor := &university.Professor{
		Name: "Cerginho",
		Code: "ABC",
	}

	code, err := storage.Set(professor)
	if err != nil {
		t.Fatal("this operation does not expect error")
	}

	if code != professor.Code {
		t.Error("unexpected Code", "code", code, "expected", professor.Code)
	}

	_, err = storage.Set(nil)
	if err == nil {
		t.Error("set method cannot accept nil values without returning an error")
	}
}

func TestProfessorStorageGet(t *testing.T) {
	storage := memory.ProfessorStorage{
		Store: make(map[string]*university.Professor),
	}

	professor := &university.Professor{
		Name: "Cerginho",
		Code: "ABC",
	}

	storage.Store["ABC"] = professor

	storedProfessor, err := storage.Get("ABC")
	if err != nil {
		t.Fatal("unexpected error", "err", err)
	}

	if storedProfessor == nil {
		t.Fatal("invalid return for storage get. professor is nil")
	}

	if professor.Name != storedProfessor.Name || professor.Code != storedProfessor.Code {
		t.Error("invalid data for stored professor")
	}
}

func TestProfessorStorageDelete(t *testing.T) {
	storage := memory.ProfessorStorage{
		Store: make(map[string]*university.Professor),
	}

	professor := &university.Professor{
		Name: "Cerginho",
		Code: "Do Ingá",
	}

	storage.Store["ABC"] = professor

	err := storage.Delete("ABC")
	if err != nil {
		t.Error("delete operation is not ok", "err", err)
	}

	if storage.Store["ABC"] != nil {
		t.Error("delete operation is not deleting from memory")
	}
}

func TestProfessorStorageGetAll(t *testing.T) {
	storage := memory.ProfessorStorage{
		Store: make(map[string]*university.Professor),
	}

	professorA := &university.Professor{
		Name: "Rogerinho",
		Code: "ABC",
	}

	storage.Store["ABC"] = professorA

	professorB := &university.Professor{
		Name: "Rogerinho",
		Code: "CBA",
	}

	storage.Store["CBA"] = professorB

	professors := []*university.Professor{
		professorA,
		professorB,
	}

	result, _ := storage.GetAll()
	for _, expected := range result {
		found := false
		for _, professor := range professors {
			if professor.Code == expected.Code {
				found = true
			}
		}

		if !found {
			t.Error("it is missing a professor in get all return")
		}
		return
	}
}

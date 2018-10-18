package memory_test

import (
	"testing"

	"github.com/alyssong/university-api/university"

	"github.com/alyssong/university-api/storage/memory"
)

func TestGet(t *testing.T) {
	classStorage := &memory.ClassStorage{
		Store: map[string]*university.Class{
			"abc": &university.Class{
				Code:      "abc",
				Professor: "321",
			},
		},
	}

	class, err := classStorage.Get("abc")
	if err != nil {
		t.Error("storage get can't throw error when item exists in the store")
	}

	if class.Code != "abc" {
		t.Error("invalid class code", "expected", "abc", "got", class.Code)
	}
}

func TestNew(t *testing.T) {
	classStorage := &memory.ClassStorage{
		Store: make(map[string]*university.Class),
	}

	err := classStorage.New("abc", "123")
	if err != nil {
		t.Error("memory storage cannot throw error when store is empty")
	}

	item, ok := classStorage.Store["abc"]
	if !ok {
		t.Error("class item is not created in map")
	}

	if item.Code != "abc" {
		t.Error("wrong class code for class creation", "expected", "abc", "got", item.Code)
	}

	if item.Professor != "123" {
		t.Error("wrong professor code for class creation", "expected", "123", "got", item.Professor)
	}

	err = classStorage.New("abc", "123")
	if err == nil {
		t.Error("memory storage should throw an error when try to create a class with a code that already exists")
	}
}

func TestAddStudents(t *testing.T) {
	classStorage := &memory.ClassStorage{
		Store: map[string]*university.Class{
			"abc": &university.Class{
				Code:      "abc",
				Professor: "321",
			},
		},
	}

	studentsCode := []string{"1", "2", "3"}
	classStorage.AddStudents("abc", studentsCode)

	for i, studentCode := range studentsCode {
		if classStorage.Store["abc"].Students[i] != studentCode {
			t.Fatal("memory storage is saving student in the wrong way")
		}
	}

	err := classStorage.AddStudents("cba", nil)
	if err == nil {
		t.Error("memory storage should throw an error when try to save students in a class that doesn't exists")
	}

	classStorage.AddStudents("abc", studentsCode)
	for i, studentCode := range classStorage.Store["abc"].Students {
		if studentCode != studentsCode[i] {
			t.Fatal("memory storage is saving student in the wrong way")
		}
	}
}

func TestRemoveStudents(t *testing.T) {
	studentsCode := []string{"1", "2", "3"}

	classStorage := &memory.ClassStorage{
		Store: map[string]*university.Class{
			"abc": &university.Class{
				Code:      "abc",
				Professor: "321",
				Students:  studentsCode,
			},
		},
	}

	classStorage.RemoveStudents("abc", []string{"1"})
	for i, code := range []string{"2", "3"} {
		if code != classStorage.Store["abc"].Students[i] {
			t.Fatal("memory storage is removing student in the wrong way")
		}
	}

	classStorage.RemoveStudents("abc", []string{"2", "3"})
	if len(classStorage.Store["abc"].Students) != 0 {
		t.Error("memory storage student slice is not empty after removing all students")
	}

	err := classStorage.RemoveStudents("cba", nil)
	if err == nil {
		t.Error("memory storage should throw an error when try to remove students from a class that doesn't exists")
	}
}

func TestSetProfessor(t *testing.T) {
	classStorage := &memory.ClassStorage{
		Store: map[string]*university.Class{
			"abc": &university.Class{
				Code:      "abc",
				Professor: "321",
			},
		},
	}

	classStorage.SetProfessor("abc", "123")
	if classStorage.Store["abc"].Professor != "123" {
		t.Error("invalid professor code after SetProfessor on memory storage")
	}

	err := classStorage.SetProfessor("cba", "")
	if err == nil {
		t.Error("memory storage should throw an error when try to set professor from a class that doesn't exists")
	}
}

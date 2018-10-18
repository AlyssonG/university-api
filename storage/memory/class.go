package memory

import (
	"github.com/alyssong/university-api/university"
	"github.com/pkg/errors"
)

type ClassStorage struct {
	Store map[string]*university.Class
}

func (cs *ClassStorage) Get(classCode string) (*university.Class, error) {
	class, ok := cs.Store[classCode]
	if !ok {
		return nil, errors.New("Class doesn't exists")
	}

	return class, nil
}

func (cs *ClassStorage) New(code, professorCode string) error {
	if _, ok := cs.Store[code]; ok {
		return errors.New("Already exists a class with this code")
	}

	class := &university.Class{
		Code:      code,
		Professor: professorCode,
	}

	cs.Store[code] = class
	return nil
}

func (cs *ClassStorage) AddStudents(classCode string, studentsCode []string) error {
	if _, ok := cs.Store[classCode]; !ok {
		return errors.New("this class doesn't exists")
	}

	for _, studentCode := range studentsCode {
		found := false
		for _, code := range cs.Store[classCode].Students {
			if code == studentCode {
				found = true
			}
		}

		if !found {
			cs.Store[classCode].Students = append(cs.Store[classCode].Students, studentCode)
		}
	}

	return nil
}

func (cs *ClassStorage) RemoveStudents(classCode string, studentsCode []string) error {
	if _, ok := cs.Store[classCode]; !ok {
		return errors.New("this class doesn't exists")
	}

	students := make([]string, 0, len(cs.Store[classCode].Students))
	for _, code := range cs.Store[classCode].Students {
		found := false
		for _, studentCode := range studentsCode {
			if code == studentCode {
				found = true
			}
		}

		if !found {
			students = append(students, code)
		}
	}

	cs.Store[classCode].Students = students
	return nil
}

func (cs *ClassStorage) SetProfessor(classCode, professorCode string) error {
	if _, ok := cs.Store[classCode]; !ok {
		return errors.New("this class doesn't exists")
	}

	cs.Store[classCode].Professor = professorCode
	return nil
}

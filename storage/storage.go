package storage

import "github.com/alyssong/university-api/university"

type Student interface {
	Set(*university.Student) (string, error)
	Get(string) (*university.Student, error)
	GetAll() ([]*university.Student, error)
	Delete(string) error
}

type Professor interface {
	Set(*university.Professor) (string, error)
	Get(string) (*university.Professor, error)
	GetAll() ([]*university.Professor, error)
	Delete(string) error
}

type Class interface {
	Get(string) (*university.Class, error)
	New(code, professorCode string) error
	AddStudents(string, []string) error
	RemoveStudents(string, []string) error
	SetProfessor(string, professorCode string) error
}

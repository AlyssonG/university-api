package storage

import "github.com/alyssong/university-api/university"

//Student defines the methods for handle with student storage.
//It defines three methods Set, Get and Delete.
//PS.:Set should be used to create or update a Student record.
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

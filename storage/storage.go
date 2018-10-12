package storage

import "github.com/alyssong/university-api/university"

//Student defines the methods for handle with student storage.
//It defines three methods Set, Get and Delete.
//PS.:Set should be used to create or update a Student record.
type Student interface {
	Set(*university.Student) (int, error)
	Get(int) (*university.Student, error)
	Delete(int) error
}

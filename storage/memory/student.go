package memory

import (
	"errors"

	"github.com/alyssong/university-api/university"
)

//StudentStorage is responsible for handle with student storage operations using memory
//as disk to store.
type StudentStorage struct {
	Store map[string]*university.Student
}

//Set creates or updates a student record in memory.
func (ss *StudentStorage) Set(student *university.Student) (string, error) {
	if student == nil {
		return "", errors.New("Student is nil. There's no data to store")
	}

	ss.Store[student.Code] = student
	return student.Code, nil
}

//Get retrives a student record in memory.
func (ss *StudentStorage) Get(studentCode string) (*university.Student, error) {
	student, ok := ss.Store[studentCode]
	if !ok {
		return nil, errors.New("Empty result")
	}

	return student, nil
}

//Delete deletes a student record in memory.
func (ss *StudentStorage) Delete(studentCode string) error {
	if ss.Store[studentCode] == nil {
		return errors.New("There is no record for this id")
	}

	delete(ss.Store, studentCode)
	return nil
}

func (ss *StudentStorage) GetAll() ([]*university.Student, error) {
	response := []*university.Student{}
	for _, value := range ss.Store {
		response = append(response, value)
	}
	return response, nil
}

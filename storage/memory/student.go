package memory

import (
	"errors"

	"github.com/alyssong/university-api/university"
)

//StudentStorage is responsible for handle with student storage operations using memory
//as disk to store.
type StudentStorage struct {
	Store map[int]*university.Student
}

//Set creates or updates a student record in memory.
func (ss *StudentStorage) Set(student *university.Student) (int, error) {
	if student == nil {
		return 0, errors.New("Student is nil. There's no data to store")
	}

	ss.Store[student.ID] = student
	return student.ID, nil
}

//Get retrives a student record in memory.
func (ss *StudentStorage) Get(studenID int) (*university.Student, error) {

	return nil, nil
}

//Delete deletes a student record in memory.
func (ss *StudentStorage) Delete(student university.Student) (int, error) {

	return 0, nil
}

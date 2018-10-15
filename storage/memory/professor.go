package memory

import (
	"errors"

	"github.com/alyssong/university-api/university"
)

//ProfessorStorage is responsible for handle with professor storage operations using memory
//as disk to store.
type ProfessorStorage struct {
	Store map[string]*university.Professor
}

//Set creates or updates a professor record in memory.
func (ss *ProfessorStorage) Set(professor *university.Professor) (string, error) {
	if professor == nil {
		return "", errors.New("Professor is nil. There's no data to store")
	}

	ss.Store[professor.Code] = professor
	return professor.Code, nil
}

//Get retrives a professor record in memory.
func (ss *ProfessorStorage) Get(professorCode string) (*university.Professor, error) {
	professor, ok := ss.Store[professorCode]
	if !ok {
		return nil, errors.New("Empty result")
	}

	return professor, nil
}

//Delete deletes a professor record in memory.
func (ss *ProfessorStorage) Delete(professorCode string) error {
	if ss.Store[professorCode] == nil {
		return errors.New("There is no record for this id")
	}

	delete(ss.Store, professorCode)
	return nil
}

func (ss *ProfessorStorage) GetAll() ([]*university.Professor, error) {
	response := []*university.Professor{}
	for _, value := range ss.Store {
		response = append(response, value)
	}
	return response, nil
}

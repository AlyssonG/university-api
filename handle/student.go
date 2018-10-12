package handle

import (
	"log"
	"net/http"

	"github.com/alyssong/university-api/storage"
)

//StudentHandle defines the API endpoints implementations for student operations
type StudentHandle struct {
	Logger  *log.Logger
	Storage storage.Student
}

//GetStudent expects a parameter ID and returns data for a student with that id
func (sh *StudentHandle) GetStudent(w http.ResponseWriter, r *http.Request) {

}

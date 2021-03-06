package handle_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alyssong/university-api/storage/memory"
	"github.com/alyssong/university-api/university"

	"github.com/alyssong/university-api/handle"
)

func TestGetStudent(t *testing.T) {
	Student := &handle.Student{
		Logger: log.New(os.Stdout, " test ", 0),
		Storage: &memory.StudentStorage{
			Store: make(map[string]*university.Student),
		},
	}

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/student?code=abc", nil)
	writter := &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}

	Student.GetStudent(writter, request)
	if writter.Code != http.StatusMethodNotAllowed {
		t.Error("get student should just accept get method")
	}

	student := &university.Student{
		Name: "Rogerinho",
		Code: "abc",
	}

	Student.Storage.Set(student)

	request = httptest.NewRequest(http.MethodGet, "http://localhost:8080/student?code=abc", nil)
	writter = &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}

	Student.GetStudent(writter, request)
	if writter.Code != http.StatusOK {
		t.Error("invalid status")
	}

	if writter.Body == nil {
		t.Fatal("nil response is not allowed")
	}

	body, err := json.Marshal(student)
	if err != nil {
		t.Fatal("error while converting student to json")
	}

	if !bytes.Equal(writter.Body.Bytes(), body) {
		t.Error("wrong return for student request")
	}

	Student.Storage.Delete("abc")
	writter = &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}

	Student.GetStudent(writter, request)
	if writter.Code != http.StatusNotFound {
		t.Error("expecting status not found when student is not in storage", "status", writter.Code)
	}

	otherStudent := &university.Student{
		Name: "Rogerinho",
		Code: "cba",
	}

	Student.Storage.Set(student)
	Student.Storage.Set(otherStudent)

	students := []*university.Student{}
	students = append(students, student, otherStudent)

	request = httptest.NewRequest(http.MethodGet, "http://localhost:8080/student", nil)
	writter = &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}
	Student.GetStudent(writter, request)

	if writter.Code != http.StatusOK {
		t.Error("invalid status for a valid request and server state")
	}

	body, err = json.Marshal(students)
	if err != nil {
		t.Fatal("error while converting student slice to json")
	}

	if !bytes.Equal(body, writter.Body.Bytes()) {
		t.Error("invalid return for api", "\ngot\n", string(writter.Body.Bytes()), "\nexpected\n", string(body))
	}
}

func TestSetStudent(t *testing.T) {
	Student := &handle.Student{
		Logger: log.New(os.Stdout, " test ", 0),
		Storage: &memory.StudentStorage{
			Store: make(map[string]*university.Student),
		},
	}

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/student", nil)
	writter := &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}

	Student.SetStudent(writter, request)
	if writter.Code != http.StatusMethodNotAllowed {
		t.Error("invalid method is accepted in set student")
	}

	student := &university.Student{
		Code: "test",
		Name: "Name test",
	}

	body, err := json.Marshal(student)
	if err != nil {
		t.Fatal("error on json marshal", "err", err)
	}

	buf := &bytes.Buffer{}
	buf.Write(body)

	request = httptest.NewRequest(http.MethodPost, "http://localhost:8080/student", buf)
	writter = &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}

	Student.SetStudent(writter, request)

	if writter.Code != http.StatusOK {
		t.Fatal("invalid status code for a valid request", "code", writter.Code)
	}

	if s, err := Student.Storage.Get("test"); s.Name != student.Name || err != nil {
		t.Error("student storage is not ok", "err", err)
	}

	body, err = json.Marshal(student)
	if err != nil {
		t.Fatal("error on json marshal", "err", err)
	}

	buf = &bytes.Buffer{}
	buf.Write(body)

	request = httptest.NewRequest(http.MethodPost, "http://localhost:8080/student", buf)

	writter = &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}

	Student.SetStudent(writter, request)

	if writter.Code != http.StatusConflict {
		t.Error("cannot create two student records with the same code", "status code", writter.Code)
	}
}

func TestDeleteStudent(t *testing.T) {
	Student := &handle.Student{
		Logger: log.New(os.Stdout, " test ", 0),
		Storage: &memory.StudentStorage{
			Store: make(map[string]*university.Student),
		},
	}

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/student/delete?code=test", nil)
	writter := &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}

	Student.DeleteStudent(writter, request)
	if writter.Code != http.StatusMethodNotAllowed {
		t.Error("delete endpoint cannot accept another http method but delete")
	}

	request = httptest.NewRequest(http.MethodDelete, "http://localhost:8080/student/delete?code=test", nil)
	writter = &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}

	Student.DeleteStudent(writter, request)
	if writter.Code != http.StatusNotFound {
		t.Error("it cant be possible to delete a non existing student record")
	}

	student := &university.Student{
		Name: "test",
		Code: "test",
	}

	Student.Storage.Set(student)

	writter = &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}

	Student.DeleteStudent(writter, request)
	if writter.Code != http.StatusOK {
		t.Error("invalid status for successful operation", "code", writter.Code)
	}
}

func TestUpdateStudent(t *testing.T) {
	Student := &handle.Student{
		Logger: log.New(os.Stdout, " test ", 0),
		Storage: &memory.StudentStorage{
			Store: make(map[string]*university.Student),
		},
	}

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/student/update?code=test", nil)
	writter := &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}

	Student.UpdateStudent(writter, request)
	if writter.Code != http.StatusMethodNotAllowed {
		t.Error("invalid http method for update student endpoint")
	}

	request = httptest.NewRequest(http.MethodPut, "http://localhost:8080/student/update?code=test", nil)
	writter = &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}

	Student.UpdateStudent(writter, request)
	if writter.Code != http.StatusNotFound {
		t.Error("updating a non existing student should throw a not found status", "got", writter.Code)
	}

	student := &university.Student{
		Name:   "test",
		Code:   "test",
		Course: "CS",
	}

	Student.Storage.Set(student)

	student.Course = "EC"

	body, _ := json.Marshal(student)
	buffer := &bytes.Buffer{}
	buffer.Write(body)

	request = httptest.NewRequest(http.MethodPut, "http://localhost:8080/student/update?code=test", buffer)
	writter = &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}

	Student.UpdateStudent(writter, request)
	if writter.Code != http.StatusOK {
		t.Fatal("invalid status code for a valid request", "got", writter.Code, "expected", http.StatusOK)
	}

	result, _ := Student.Storage.Get(student.Code)
	if result.Course != student.Course {
		t.Error("update student endpoint is not updating record in storage", "expected", student.Course, "got", result.Course)
	}
}

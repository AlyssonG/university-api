package handle_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/alyssong/university-api/storage/memory"
	"github.com/alyssong/university-api/university"

	"github.com/alyssong/university-api/handle"
)

func TestGetStudent(t *testing.T) {
	studentHandle := &handle.StudentHandle{
		Logger: log.New(os.Stdout, " test ", 0),
		Storage: &memory.StudentStorage{
			Store: make(map[string]*university.Student),
		},
	}

	student := &university.Student{
		Name: "Rogerinho",
		Code: "abc",
	}

	studentHandle.Storage.Set(student)
	writter := &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}

	request := &http.Request{
		URL: &url.URL{
			RawQuery: "id=abc",
		},
	}

	studentHandle.GetStudent(writter, request)
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

	studentHandle.Storage.Delete("abc")
	writter = &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}

	studentHandle.GetStudent(writter, request)
	if writter.Code != http.StatusNotFound {
		t.Error("expecting status not found when student is not in storage", "status", writter.Code)
	}
}

func TestSetStudent(t *testing.T) {
	studentHandle := &handle.StudentHandle{
		Logger: log.New(os.Stdout, " test ", 0),
		Storage: &memory.StudentStorage{
			Store: make(map[string]*university.Student),
		},
	}

	writter := &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
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

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/student", buf)

	studentHandle.SetStudent(writter, request)

	if writter.Code != http.StatusOK {
		t.Fatal("invalid status code for a valid request", "code", writter.Code)
	}

	if s, err := studentHandle.Storage.Get("test"); s.Name != student.Name || err != nil {
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

	studentHandle.SetStudent(writter, request)

	if writter.Code != http.StatusConflict {
		t.Error("cannot create two student records with the same code", "status code", writter.Code)
	}
}

func TestDeleteStudent(t *testing.T) {
	studentHandle := &handle.StudentHandle{
		Logger: log.New(os.Stdout, " test ", 0),
		Storage: &memory.StudentStorage{
			Store: make(map[string]*university.Student),
		},
	}

	writter := &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/student/delete?code=test", nil)

	studentHandle.DeleteStudent(writter, request)
	if writter.Code != http.StatusNotFound {
		t.Error("it cant be possible to delete a non existing student record")
	}

	student := &university.Student{
		Name: "test",
		Code: "test",
	}

	studentHandle.Storage.Set(student)

	writter = &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}

	studentHandle.DeleteStudent(writter, request)
	if writter.Code != http.StatusOK {
		t.Error("invalid status for successful operation", "code", writter.Code)
	}
}

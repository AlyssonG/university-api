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
			Store: make(map[int]*university.Student),
		},
	}

	student := &university.Student{
		ID:   1,
		Name: "Rogerinho",
		Code: "Do Ingá",
	}

	studentHandle.Storage.Set(student)
	writter := &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}

	request := &http.Request{
		URL: &url.URL{
			RawQuery: "id=1",
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

	studentHandle.Storage.Delete(1)

	writter = &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}

	studentHandle.GetStudent(writter, request)
	if writter.Code != http.StatusNotFound {
		t.Error("expecting status not found when student is not in storage", "status", writter.Code)
	}

	writter = &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}

	request = &http.Request{
		URL: &url.URL{
			RawQuery: "id=abc",
		},
	}

	studentHandle.GetStudent(writter, request)
	if writter.Code != http.StatusInternalServerError {
		t.Error("expecting status internal server error when student id is invalid", "status", writter.Code)
	}

	writter = &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}

	request = &http.Request{
		URL: &url.URL{
			RawQuery: "",
		},
	}

	studentHandle.GetStudent(writter, request)
	if writter.Code != http.StatusInternalServerError {
		t.Error("expecting status internal server error when student id is invalid", "status", writter.Code)
	}

}

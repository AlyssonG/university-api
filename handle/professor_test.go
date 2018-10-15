package handle_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alyssong/university-api/handle"
	"github.com/alyssong/university-api/storage/memory"
	"github.com/alyssong/university-api/university"
)

func setup() *handle.Professor {
	return &handle.Professor{
		Logger: log.New(os.Stdout, "-- Test --", log.LstdFlags),
		Storage: &memory.ProfessorStorage{
			Store: map[string]*university.Professor{},
		},
	}
}

func buildRequest(method, url string, body io.Reader) (*httptest.ResponseRecorder, *http.Request) {
	request := httptest.NewRequest(method, url, body)
	writter := &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}

	return writter, request
}

func bodyBuilder(data interface{}) []byte {
	bytes, err := json.Marshal(data)
	if err != nil {
		return []byte("{}")
	}
	return bytes
}

func TestGetProfessor(t *testing.T) {
	professorHandle := setup()

	w, r := buildRequest(http.MethodPost, "http://localhost:8080", nil)
	professorHandle.GetProfessor(w, r)

	if w.Code != http.StatusMethodNotAllowed {
		t.Error("accepted wrong http method. should return", http.StatusMethodNotAllowed, "but got", w.Code)
	}

	tctc := []struct {
		Name, URL    string
		ExpectedBody []byte
		ExpectedCode int
		Professors   []*university.Professor
	}{
		{
			Name:         "with empty storage",
			URL:          "professor?code=1",
			ExpectedCode: http.StatusNotFound,
		},
		{
			Name:         "with one professor",
			URL:          "professor?code=1",
			ExpectedCode: http.StatusOK,
			Professors: []*university.Professor{
				{Name: "test", Code: "1"},
			},
			ExpectedBody: bodyBuilder(&university.Professor{Name: "test", Code: "1"}),
		},
		{
			Name:         "with two professor",
			URL:          "professor",
			ExpectedCode: http.StatusOK,
			Professors: []*university.Professor{
				{Name: "test", Code: "1"},
				{Name: "test2", Code: "2"},
			},
			ExpectedBody: bodyBuilder([]*university.Professor{
				{Name: "test", Code: "1"},
				{Name: "test2", Code: "2"}},
			),
		},
	}

	for _, tc := range tctc {
		t.Run(tc.Name, func(t *testing.T) {
			professorHandle := setup()
			for _, professor := range tc.Professors {
				professorHandle.Storage.Set(professor)
			}

			w, r := buildRequest(http.MethodGet, fmt.Sprintf("http://localhost:8080/%s", tc.URL), nil)
			professorHandle.GetProfessor(w, r)

			if w.Code != tc.ExpectedCode {
				t.Fatal("unexpected code for", tc.Name, "test.", "expected", tc.ExpectedCode, "got", w.Code)
			}

			if tc.ExpectedBody == nil {
				return
			}

			if !bytes.Equal(tc.ExpectedBody, w.Body.Bytes()) {
				t.Error("wrong body on", tc.Name, "test", "expected", string(tc.ExpectedBody), "got", string(w.Body.Bytes()))
			}
		})
	}

}
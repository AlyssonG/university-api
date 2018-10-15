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

func TestSetProfessor(t *testing.T) {
	professorHandle := setup()

	w, r := buildRequest(http.MethodGet, "http://localhost:8080", nil)
	professorHandle.SetProfessor(w, r)

	if w.Code != http.StatusMethodNotAllowed {
		t.Error("accepting wrong http method. should return", http.StatusMethodNotAllowed, "but got", w.Code)
	}

	tctc := []struct {
		Name                 string
		StatusCode           int
		Professors           []*university.Professor
		ProfessorCode        string
		ExpectedResult, Body []byte
	}{
		{
			Name:       "with existing code",
			Body:       bodyBuilder(&university.Professor{Code: "1"}),
			StatusCode: http.StatusConflict,
			Professors: []*university.Professor{
				{
					Code: "1",
				},
			},
		},
		{
			Name:           "set one professor",
			Body:           bodyBuilder(&university.Professor{Code: "1"}),
			StatusCode:     http.StatusOK,
			ProfessorCode:  "1",
			ExpectedResult: bodyBuilder(&university.Professor{Code: "1"}),
		},
	}

	for _, tc := range tctc {
		t.Run(tc.Name, func(t *testing.T) {
			professorHandle := setup()
			for _, professor := range tc.Professors {
				professorHandle.Storage.Set(professor)
			}

			buffer := &bytes.Buffer{}
			buffer.Write(tc.Body)

			w, r := buildRequest(http.MethodPost, "http://localhost:8080", buffer)
			professorHandle.SetProfessor(w, r)

			if w.Code != tc.StatusCode {
				t.Error("unexpected status.", "expected", tc.StatusCode, "but got", w.Code)
			}

			if tc.ExpectedResult == nil {
				return
			}

			professor, _ := professorHandle.Storage.Get(tc.ProfessorCode)
			result := bodyBuilder(professor)

			if !bytes.Equal(result, tc.ExpectedResult) {
				t.Error("unexpected result for set professor", "expected", string(tc.ExpectedResult), "got", string(result))
			}
		})
	}
}

func TestDeleteProfessor(t *testing.T) {
	professorHandle := setup()

	w, r := buildRequest(http.MethodGet, "http://localhost:8000", nil)
	professorHandle.DeleteProfessor(w, r)

	if w.Code != http.StatusMethodNotAllowed {
		t.Error("accepting wrong http method. should return", http.StatusMethodNotAllowed, "but got", w.Code)
	}

	tctc := []struct {
		Name       string
		Professors []*university.Professor
		Code       string
		URL        string
		StatusCode int
	}{
		{
			Name:       "delete with empty storage",
			URL:        "code=1",
			StatusCode: http.StatusNotFound,
		},
		{
			Name:       "delete one professor",
			URL:        "code=1",
			Code:       "1",
			StatusCode: http.StatusOK,
			Professors: []*university.Professor{
				{
					Code: "1",
				},
			},
		},
	}

	for _, tc := range tctc {
		professorHandle := setup()
		for _, professor := range tc.Professors {
			professorHandle.Storage.Set(professor)
		}

		w, r := buildRequest(http.MethodDelete, fmt.Sprintf("http://localhost:8000/professor/delete?%s", tc.URL), nil)
		professorHandle.DeleteProfessor(w, r)

		if w.Code != tc.StatusCode {
			t.Error("unexpected status code. should return", tc.StatusCode, "but got", w.Code)
		}

	}
}

func TestUpdateProfessor(t *testing.T) {
	professorHandle := setup()

	w, r := buildRequest(http.MethodGet, "http://localhost:8080", nil)
	professorHandle.UpdateProfessor(w, r)

	if w.Code != http.StatusMethodNotAllowed {
		t.Error("accepting wrong http method. should return", http.StatusMethodNotAllowed, "but got", w.Code)
	}

	tctc := []struct {
		Name                 string
		StatusCode           int
		Professors           []*university.Professor
		ProfessorCode        string
		ExpectedResult, Body []byte
	}{
		{
			Name:       "with empty storage",
			Body:       bodyBuilder(&university.Professor{Code: "1"}),
			StatusCode: http.StatusNotFound,
		},
		{
			Name:           "set one professor",
			Body:           bodyBuilder(&university.Professor{Name: "changed name"}),
			StatusCode:     http.StatusOK,
			Professors:     []*university.Professor{{Name: "test", Code: "1"}},
			ProfessorCode:  "1",
			ExpectedResult: bodyBuilder(&university.Professor{Code: "1", Name: "changed name"}),
		},
	}

	for _, tc := range tctc {
		t.Run(tc.Name, func(t *testing.T) {
			professorHandle := setup()
			for _, professor := range tc.Professors {
				professorHandle.Storage.Set(professor)
			}

			buffer := &bytes.Buffer{}
			buffer.Write(tc.Body)

			w, r := buildRequest(http.MethodPut, fmt.Sprintf("http://localhost:8080?code=%s", tc.ProfessorCode), buffer)
			professorHandle.UpdateProfessor(w, r)

			if w.Code != tc.StatusCode {
				t.Error("unexpected status.", "expected", tc.StatusCode, "but got", w.Code)
			}

			if tc.ExpectedResult == nil {
				return
			}

			professor, _ := professorHandle.Storage.Get(tc.ProfessorCode)
			result := bodyBuilder(professor)

			if !bytes.Equal(result, tc.ExpectedResult) {
				t.Error("unexpected result for set professor", "expected", string(tc.ExpectedResult), "got", string(result))
			}
		})
	}
}

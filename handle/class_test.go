package handle_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alyssong/university-api/university"

	"github.com/alyssong/university-api/handle"
	"github.com/alyssong/university-api/storage/memory"
)

//TODO: Continuar testes. Dá para fazer olhando os testes de student e professor

func TestSetClass(t *testing.T) {
	classHandle := &handle.Class{
		Logger: log.New(os.Stdout, "", log.LstdFlags),
		Storage: &memory.ClassStorage{
			Store: make(map[string]*university.Class),
		},
	}

	request := httptest.NewRequest(http.MethodGet, "http://", nil)
	response := &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}

	classHandle.SetClass(response, request)
	if response.Code != http.StatusMethodNotAllowed {
		t.Error("codigo de retorno invalido", "esperando", http.StatusMethodNotAllowed)
	}

	class := university.Class{
		Professor: "123",
		Code:      "A1",
	}

	bytesBody, _ := json.Marshal(class)

	body := &bytes.Buffer{}
	body.Write(bytesBody)

	request = httptest.NewRequest(http.MethodPost, "http://", body)
	response = &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}

	classHandle.SetClass(response, request)
	if response.Code != http.StatusOK {
		t.Error("status deve ser ok mas retornou", response.Code)
		t.Error(string(response.Body.Bytes()))
	}

	if class, _ := classHandle.Storage.Get("A1"); class == nil {
		t.Error("turma nao criada")
	}

	body = &bytes.Buffer{}
	body.Write(bytesBody)

	request = httptest.NewRequest(http.MethodPost, "http://", body)
	response = &httptest.ResponseRecorder{
		Body: &bytes.Buffer{},
	}

	classHandle.SetClass(response, request)
	if response.Code != http.StatusConflict {
		t.Error("tentando criar uma turma que ja existe deve dar conflito", "mas retornou", response.Code)
	}
}

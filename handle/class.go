package handle

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/alyssong/university-api/university"

	"github.com/alyssong/university-api/storage"
)

//TODO: Criar testes para AddStudents, SetProfessor e RemoveStudents
type Class struct {
	Logger  *log.Logger
	Storage storage.Class
}

func (c *Class) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.GetClass(w, r)
	case http.MethodPost:
		c.SetClass(w, r)
	case http.MethodPut:
		c.AddStudents(w, r)
	case http.MethodPatch:
		c.SetProfessor(w, r)
	case http.MethodDelete:
		c.RemoveStudents(w, r)
	default:
		http.Error(w, "{\"message\": \"invalid method\"}", http.StatusMethodNotAllowed)
	}
}

func (c *Class) SetClass(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, "{\"message\": \"error while reading request body\"}", http.StatusInternalServerError)
		c.Logger.Println("error while reading request body", "err", err)
		return
	}

	class := &university.Class{}
	json.Unmarshal(body, class)

	exists, err := c.Storage.Get(class.Code)
	if exists != nil {
		http.Error(w, "", http.StatusConflict)
		return
	}

	c.Storage.New(class.Code, class.Professor)
	w.Write([]byte("{\"status\": \"ok\"}"))
}

func (c *Class) GetClass(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	code := query.Get("code")

	class, err := c.Storage.Get(code)
	if err != nil {
		http.Error(w, "{\"message\": \"error while getting class\"}", http.StatusInternalServerError)
		c.Logger.Println("error while getting class", "err", err)
		return
	}

	bytes, err := json.Marshal(class)
	if err != nil {
		http.Error(w, "{\"message\": \"error in response format\"}", http.StatusInternalServerError)
		c.Logger.Println("error on struct conversion to json", "err", err)
		return
	}

	w.Write(bytes)
}

func (c *Class) AddStudents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	code := query.Get("code")

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, "{\"message\": \"error while reading request body\"}", http.StatusInternalServerError)
		c.Logger.Println("error while reading request body", "err", err)
		return
	}

	class := &university.Class{}
	err = json.Unmarshal(body, class)
	if err != nil {
		http.Error(w, "{\"message\": \"invalid request body\"}", http.StatusInternalServerError)
		c.Logger.Println("invalid request body", "err", err)
		return
	}

	err = c.Storage.AddStudents(code, class.Students)
	if err != nil {
		http.Error(w, "{\"message\": \"error on add students operation\"}", http.StatusInternalServerError)
		c.Logger.Println("error on add students operation", "err", err)
		return
	}

	w.Write([]byte("{\"status\": \"ok\"}"))
}

func (c *Class) RemoveStudents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	code := query.Get("code")

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, "{\"message\": \"error while reading request body\"}", http.StatusInternalServerError)
		c.Logger.Println("error while reading request body", "err", err)
		return
	}

	class := &university.Class{}
	err = json.Unmarshal(body, class)
	if err != nil {
		http.Error(w, "{\"message\": \"invalid request body\"}", http.StatusInternalServerError)
		c.Logger.Println("invalid request body", "err", err)
		return
	}

	err = c.Storage.RemoveStudents(code, class.Students)
	if err != nil {
		http.Error(w, "{\"message\": \"error on remove students operation\"}", http.StatusInternalServerError)
		c.Logger.Println("error on remove students operation", "err", err)
		return
	}

	w.Write([]byte("{\"status\": \"ok\"}"))
}

func (c *Class) SetProfessor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	code := query.Get("code")

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, "{\"message\": \"error while reading request body\"}", http.StatusInternalServerError)
		c.Logger.Println("error while reading request body", "err", err)
		return
	}

	class := &university.Class{}
	err = json.Unmarshal(body, class)
	if err != nil {
		http.Error(w, "{\"message\": \"invalid request body\"}", http.StatusInternalServerError)
		c.Logger.Println("invalid request body", "err", err)
		return
	}

	//TODO: Checar se o professor existe
	err = c.Storage.SetProfessor(code, class.Professor)
	if err != nil {
		http.Error(w, "{\"message\": \"error on set professor operation\"}", http.StatusInternalServerError)
		c.Logger.Println("error on set professor operation", "err", err)
		return
	}

	w.Write([]byte("{\"status\": \"ok\"}"))
}

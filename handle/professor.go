package handle

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sort"

	"github.com/alyssong/university-api/university"

	"github.com/alyssong/university-api/storage"
)

type Professor struct {
	Logger  *log.Logger
	Storage storage.Professor
}

func (p *Professor) GetProfessor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "{}", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	code := query.Get("code")

	if code != "" {
		p.getSpecificProfessor(code, w, r)
	} else {
		p.getAllProfessors(w, r)
	}

}

func (p *Professor) getSpecificProfessor(professorCode string, w http.ResponseWriter, r *http.Request) {
	professor, err := p.Storage.Get(professorCode)
	if err != nil {
		http.Error(w, "{\"message\": \"error while recovering professor from storage\"}", http.StatusNotFound)
		p.Logger.Println("error while recovering professor from storage", "err", err)
		return
	}

	response, err := json.Marshal(professor)
	if err != nil {
		http.Error(w, "{\"message\": \"error while generating professor json\"}", http.StatusInternalServerError)
		p.Logger.Println("error while generating professor json", "err", err)
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	w.Write(response)
}

func (p *Professor) getAllProfessors(w http.ResponseWriter, r *http.Request) {
	professors, err := p.Storage.GetAll()
	if err != nil {
		http.Error(w, "{\"message\": \"error while recovering professors from storage\"}", http.StatusNotFound)
		p.Logger.Println("error while recovering professors from storage", "err", err)
		return
	}

	sort.Slice(professors, func(i, j int) bool { return professors[i].Code < professors[j].Code })

	response, err := json.Marshal(professors)
	if err != nil {
		http.Error(w, "{\"message\": \"error while generating professors json\"}", http.StatusInternalServerError)
		p.Logger.Println("error while generating professors json", "err", err)
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	w.Write(response)
}

func (p *Professor) SetProfessor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "{]", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "{\"message\": \"error while reading body\"}", http.StatusInternalServerError)
		p.Logger.Println("error while reading body", "err", err)
		return
	}

	professor := &university.Professor{}
	err = json.Unmarshal(body, professor)
	if err != nil {
		http.Error(w, "{\"message\": \"error while parsing body\"}", http.StatusInternalServerError)
		p.Logger.Println("error while parsing body", "err", err)
		return
	}

	storedProfessor, err := p.Storage.Get(professor.Code)
	if storedProfessor != nil {
		http.Error(w, "{\"message\": \"professor already exists with this code\"}", http.StatusConflict)
		return
	}

	_, err = p.Storage.Set(professor)
	if err != nil {
		http.Error(w, "{\"message\": \"error while creating professor record on db\"}", http.StatusInternalServerError)
		p.Logger.Println("error while creating professor record on db", "err", err)
		return
	}

	w.Write([]byte("{\"status\": \"success\"}"))
}

func (p *Professor) DeleteProfessor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "{}", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	code := query.Get("code")

	professor, err := p.Storage.Get(code)
	if professor == nil || err != nil {
		http.Error(w, "{\"message\": \"cannot find professor to delete\"}", http.StatusNotFound)
		p.Logger.Println("error while finding professor to delete", "err", err)
		return
	}

	err = p.Storage.Delete(code)
	if err != nil {
		http.Error(w, "{\"message\": \"error while deleting professor\"}", http.StatusInternalServerError)
		p.Logger.Println("error while deleting professor", "err", err)
		return
	}

	w.Write([]byte("{\"status\": \"success\"}"))
}

func (p *Professor) UpdateProfessor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "{}", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	code := query.Get("code")

	professor, _ := p.Storage.Get(code)
	if professor == nil {
		http.Error(w, "{\"message\": \"professor not found\"}", http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "{\"message\": \"error while reading request body\"}", http.StatusInternalServerError)
		p.Logger.Println("error while reading request body", "err", err)
	}

	professorData := &university.Professor{}
	err = json.Unmarshal(body, professorData)
	if err != nil {
		http.Error(w, "{\"message\": \"error while parsing body\"}", http.StatusInternalServerError)
		p.Logger.Println("error while parsing professor body to struct", "err", err)
		return
	}

	professor.Name = professorData.Name
	professor.Department = professorData.Department

	_, err = p.Storage.Set(professor)
	if err != nil {
		http.Error(w, "{\"message\": \"error on professor storage\"}", http.StatusInternalServerError)
		p.Logger.Println("cannot update professor record in storage", "err", err)
		return
	}

	w.Write([]byte("{\"status\": \"success\"}"))
}

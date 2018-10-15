package handle

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"

	"github.com/alyssong/university-api/storage"
)

type Professor struct {
	Logger  *log.Logger
	Storage storage.Professor
}

func (p *Professor) GetProfessor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusMethodNotAllowed)
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
		http.Error(w, "error while recovering professor from storage", http.StatusNotFound)
		p.Logger.Println("error while recovering professor from storage", "err", err)
		return
	}

	response, err := json.Marshal(professor)
	if err != nil {
		http.Error(w, "error while generating professor json", http.StatusInternalServerError)
		p.Logger.Println("error while generating professor json", "err", err)
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	w.Write(response)
}

func (p *Professor) getAllProfessors(w http.ResponseWriter, r *http.Request) {
	professors, err := p.Storage.GetAll()
	if err != nil {
		http.Error(w, "error while recovering professors from storage", http.StatusNotFound)
		p.Logger.Println("error while recovering professors from storage", "err", err)
		return
	}

	sort.Slice(professors, func(i, j int) bool { return professors[i].Code < professors[j].Code })

	response, err := json.Marshal(professors)
	if err != nil {
		http.Error(w, "error while generating professors json", http.StatusInternalServerError)
		p.Logger.Println("error while generating professors json", "err", err)
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	w.Write(response)
}

package handle

import (
	"encoding/json"
	"log"
	"net/http"

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

	professor, err := p.Storage.Get(code)
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

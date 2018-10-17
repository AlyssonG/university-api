package university

type Class struct {
	Code      string   `json:"code"`
	Professor string   `json:"professor"`
	Students  []string `json:"students"`
}

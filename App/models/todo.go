package models

type Todo struct {
	Id        string `json:"_key"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type Todos struct {
	Todos []Todo `json:"todos"`
}

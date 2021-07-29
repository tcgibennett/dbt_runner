package models

type ResponseCompiledSQL struct {
	Name string `json:"name"`
	SQL  string `json:"sql"`
}

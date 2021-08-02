package models

type ResponseCompiledSQL struct {
	Name string `json:"name"`
	SQL  string `json:"sql"`
}

type DBTProject struct {
	Name string `yaml:"name"`
	Vars map[string]string `yaml:"vars"`
}

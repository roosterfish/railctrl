package api

type Route struct {
	Path []Track `yaml:"path" json:"path"`
	Cost int     `yaml:"cost" json:"cost"`
}

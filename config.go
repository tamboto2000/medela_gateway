package medelagateway

import (
	"encoding/json"
	"os"
)

type Config struct {
	Endpoints []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	Endpoint    string       `json:"endpoint"`
	Method      string       `json:"method"`
	Backend     Backend      `json:"backend"`
	Middlewares []Middleware `json:"middlewares"`
}

type Backend struct {
	Host       string `json:"host"`
	UrlPattern string `json:"url_pattern"`
	Method     string `json:"method"`
}

type Middleware struct {
	Host                string `json:"host"`
	UrlPattern          string `json:"url_pattern"`
	MergeResponseBody   bool   `json:"merge_response_body"`
	MergeResponseHeader bool   `json:"merge_response_header"`
}

func ParseConfigFromFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	conf := new(Config)
	if err := json.NewDecoder(f).Decode(conf); err != nil {
		return nil, err
	}

	return conf, nil
}

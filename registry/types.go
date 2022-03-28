package registry

import (
	"github.com/getkin/kin-openapi/openapi3"
	"tpm-symphony/registry/crawler"
)

type Config struct {
	CrawlerCfg crawler.Config `yaml:"crawler" mapstructure:"crawler" json:"crawler"`
}

type SymphonyPathExtension struct {
	OpPath    string `yaml:"op-path" mapstructure:"op-path" json:"op-path"`
	GroupPath string `yaml:"group-path" mapstructure:"group-path" json:"group-path" `
}

type SymphonyOperationExtension struct {
	HttpMethod string `yaml:"http-method" mapstructure:"http-method" json:"http-method"`
	Id         string `yaml:"id" mapstructure:"id" json:"id"`
}

func (o SymphonyOperationExtension) IsZero() bool {
	return o.Id == ""
}

type Resource struct {
	Name       string `yaml:"name" mapstructure:"name" json:"name"`
	Method     string `yaml:"method" mapstructure:"method" json:"method"`
	Path       string `yaml:"path" mapstructure:"path" json:"path"`
	SymphonyId string `yaml:"sid" mapstructure:"sid" json:"sid"`
}

type ResourceGroup struct {
	Name string `yaml:"name" mapstructure:"name" json:"name"`
	Path string `yaml:"path" mapstructure:"path" json:"path" `
}

type Entry struct {
	Repo          crawler.OrchestrationRepo
	OpenapiDoc    *openapi3.T
	ResourceGroup ResourceGroup `yaml:"resource-group" mapstructure:"resource-group" json:"resource-group"`
	Resources     []Resource    `yaml:"resources" mapstructure:"resources" json:"resources"`
}

type OrchestrationRegistry struct {
	Entries []Entry `yaml:"entries" mapstructure:"entries" json:"entries"`
}

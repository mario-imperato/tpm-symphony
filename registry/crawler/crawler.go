package crawler

import (
	"io/ioutil"
	"os"
	"path"
	"regexp"
)

type Config struct {
	Type       string `yaml:"type" mapstructure:"type" json:"type"`
	MountPoint string `yaml:"mount-point" mapstructure:"mount-point" json:"mount-point"`
}

type Asset struct {
	Type string `yaml:"type" json:"type" mapstructure:"type"`
	Path string `yaml:"path" json:"path" mapstructure:"path"`
	Data []byte `yaml:"data" json:"data" mapstructure:"data"`
}

func (a Asset) IsZero() bool {
	return a.Type == ""
}

type EndPointAsset struct {
	Id string `yaml:"id" mapstructure:"id"`
}

type OrchestrationAsset struct {
	ExecutionGraph Asset           `yaml:"exec-graph" mapstructure:"exec-graph"`
	Endpoints      []EndPointAsset `yaml:"endpoints" mapstructure:"endpoints"`
}

func Crawl(cfg *Config) (Asset, []OrchestrationAsset, error) {

	fi, err := scan4OpenApiFiles(cfg.MountPoint)
	if err != nil {
		return Asset{}, nil, err
	}

	openApiPath := path.Join(cfg.MountPoint, fi.Name())
	return Asset{Type: "open-api", Path: openApiPath}, nil, nil
}

var OpenapiFileRegexp = regexp.MustCompile("tpm-symphony-openapi\\.(yml|yaml)")

func scan4OpenApiFiles(dir string) (os.FileInfo, error) {

	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, fi := range fis {
		if !fi.IsDir() {
			if OpenapiFileRegexp.Match([]byte(fi.Name())) {
				return fi, nil
			}
		}
	}

	return nil, nil
}

package crawler

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"path"
	"tpm-symphony/constants"
)

type Asset struct {
	Name string  `yaml:"name" json:"name" mapstructure:"name"`
	Type string  `yaml:"type" json:"type" mapstructure:"type"`
	Path string  `yaml:"path" json:"path" mapstructure:"path"`
	Data []byte  `yaml:"data" json:"data" mapstructure:"data"`
	Refs []Asset `yaml:"references" json:"data" mapstructure:"references"`
}

func (a Asset) IsZero() bool {
	return a.Type == ""
}

type OrchestrationRepo struct {
	path           string             `yaml:"path" mapstructure:"path"`
	apiDefinition  Asset              `yaml:"api" mapstructure:"api"`
	assets         []Asset            `yaml:"assets" mapstructure:"assets"`
	orchestrations map[string][]Asset `yaml:"orchestrations" mapstructure:"orchestrations"`
}

func (r *OrchestrationRepo) ShowInfo() {
	log.Info().Str(constants.SemLogPath, r.GetPath()).Int("num-assets", len(r.assets)).Msg("BOF repo information --------")
	log.Info().Str(constants.SemLogPath, r.GetPath()).Str(constants.SemLogFile, r.apiDefinition.Path).Msg(constants.SemLogOpenApi)
	for _, a := range r.assets {
		log.Info().Str(constants.SemLogType, a.Type).Str(constants.SemLogPath, r.GetPath()).Str(constants.SemLogFile, a.Path).Msg("open-api external value")
	}

	for no, o := range r.orchestrations {
		log.Info().Str(constants.SemLogOrchestrationSid, no).Msg("orchestration")
		for _, a := range o {
			log.Info().Str(constants.SemLogType, a.Type).Str(constants.SemLogPath, r.GetPath()).Str(constants.SemLogFile, a.Path).Msg("orchestration asset")
		}
	}

}

func (r *OrchestrationRepo) GetPath() string {
	return r.path
}

func (r *OrchestrationRepo) GetOpenApiData() (string, []byte, error) {

	if r.apiDefinition.IsZero() {
		return "", nil, fmt.Errorf("no api definition found in repo %s", r.path)
	}

	b, err := r.getData(r.apiDefinition.Path)
	return r.apiDefinition.Path, b, err
}

func (r *OrchestrationRepo) getData(fn string) ([]byte, error) {

	if r.apiDefinition.IsZero() {
		return nil, fmt.Errorf("no api definition found in repo %s", r.path)
	}

	if r.apiDefinition.Data != nil {
		return r.apiDefinition.Data, nil
	}

	resolvedPath := path.Join(r.path, r.apiDefinition.Path)
	b, err := ioutil.ReadFile(resolvedPath)
	if err != nil {
		return nil, err
	}

	r.apiDefinition.Data = b
	return r.apiDefinition.Data, nil
}

func (r *OrchestrationRepo) GetOpenApiAssets(fn string) ([]byte, error) {

	if len(r.assets) == 0 {
		return nil, fmt.Errorf("no assets references found in repo %s", r.path)
	}

	ndx := findAssetIndexByPath(r.assets, fn)
	if ndx < 0 {
		return nil, fmt.Errorf("no assets references found in repo %s for %s", r.path, fn)
	}

	if r.assets[ndx].Data != nil {
		return r.assets[ndx].Data, nil
	}

	resolvedPath := path.Join(r.path, r.assets[ndx].Path)
	b, err := ioutil.ReadFile(resolvedPath)
	if err != nil {
		return nil, err
	}

	r.assets[ndx].Data = b
	return r.assets[ndx].Data, nil
}

func findAssetIndexByPath(assets []Asset, p string) int {
	for i, a := range assets {
		if a.Path == p {
			return i
		}
	}

	return -1
}

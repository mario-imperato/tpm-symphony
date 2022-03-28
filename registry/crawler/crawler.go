package crawler

import (
	"github.com/mario-imperato/tpm-common/util"
	"path"
)

type Config struct {
	Type       string `yaml:"type" mapstructure:"type" json:"type"`
	MountPoint string `yaml:"mount-point" mapstructure:"mount-point" json:"mount-point"`
}

var DefaultIgnoreList = []string{
	"^\\.",
	"tpm-symphony-openapi\\.(yml|yaml)",
}

var OpenApiFileFindIncludeList = []string{
	"tpm-symphony-openapi\\.(yml|yaml)",
}

func Crawl(cfg *Config) ([]OrchestrationRepo, error) {

	fi, err := scan4OpenApiFiles(cfg.MountPoint)
	if err != nil {
		return nil, err
	}

	repos := make([]OrchestrationRepo, 0)
	for _, f := range fi {
		openApiPath := f
		repo := OrchestrationRepo{
			path:          path.Dir(openApiPath),
			apiDefinition: Asset{Name: path.Base(openApiPath), Type: "open-api", Path: path.Base(openApiPath)},
		}

		assets, err := scan4AssetFiles(repo.path)
		if err != nil {
			return nil, err
		}
		repo.assets = assets

		subFolders, err := util.FindFiles(repo.path, util.WithFindFileType(util.FileTypeDir), util.WithFindOptionIgnoreList(DefaultIgnoreList))
		if err != nil {
			return nil, err
		}

		m := make(map[string][]Asset, 0)
		for _, fld := range subFolders {
			fldName, assets, err := scan4OrchestrationFiles(fld)
			if err != nil {
				return nil, err
			}

			m[fldName] = assets
		}

		repo.orchestrations = m
		repos = append(repos, repo)
	}

	return repos, nil
}

func scan4OpenApiFiles(dir string) ([]string, error) {
	files, err := util.FindFiles(dir, util.WithFindOptionNavigateSubDirs(), util.WithFindOptionIncludeList(OpenApiFileFindIncludeList))
	if err != nil {
		return nil, err
	}

	return files, nil
}

func scan4AssetFiles(dir string) ([]Asset, error) {
	files, err := util.FindFiles(dir, util.WithFindFileType(util.FileTypeFile), util.WithFindOptionIgnoreList(DefaultIgnoreList))
	if err != nil {
		return nil, err
	}

	assets := make([]Asset, 0, len(files))
	for _, a := range files {
		assets = append(assets, Asset{Name: path.Base(a), Type: "else", Path: path.Base(a)})
	}

	return assets, nil
}

func scan4OrchestrationFiles(dir string) (string, []Asset, error) {
	files, err := util.FindFiles(dir, util.WithFindFileType(util.FileTypeFile), util.WithFindOptionIgnoreList(DefaultIgnoreList))
	if err != nil {
		return "", nil, err
	}

	assets := make([]Asset, 0, len(files))
	for _, a := range files {
		assets = append(assets, Asset{Type: "else", Path: path.Join(path.Base(dir), path.Base(a))})
	}

	return path.Base(dir), assets, nil
}

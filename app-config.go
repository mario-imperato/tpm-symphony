package main

import (
	_ "embed"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-gin/httpsrv"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-gin/middleware"
	"github.com/rs/zerolog/log"
	"gitlab.alm.poste.it/go/configuration"
	"os"
	"strings"
	"tpm-symphony/registry"
)

type AppConfig struct {
	Registry   registry.Config                     `yaml:"registry" mapstructure:"registry" json:"registry"`
	Http       httpsrv.Config                      `yaml:"http" mapstructure:"http" json:"http"`
	MwRegistry *middleware.MwHandlerRegistryConfig `yaml:"mw-handler-registry" mapstructure:"mw-handler-registry" json:"mw-handler-registry"`
}

// Default config file.
//go:embed config.yml
var projectConfigFile []byte

const ConfigFileEnvVar = "SYMPHONY_CFG_FILE_PATH"
const ConfigurationName = "tpm-symphony"

func ReadConfig() (*AppConfig, error) {

	configPath := os.Getenv(ConfigFileEnvVar)
	var cfgFileReader *strings.Reader
	if configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			log.Info().Str("cfg-file-name", configPath).Msg("reading config")
			cfgContent, rerr := util.ResolveEnvVarsInConfigFile(configPath)
			if rerr != nil {
				return nil, err
			} else {
				cfgFileReader = strings.NewReader(cfgContent)
			}

		} else {
			return nil, fmt.Errorf("the %s env variable has been set but no file cannot be found at %s", ConfigFileEnvVar, configPath)
		}
	} else {
		log.Warn().Msgf("The config path variable %s has not been set. Reverting to bundled configuration", ConfigFileEnvVar)
		cfgFileReader = strings.NewReader(string(projectConfigFile))

		// return nil, fmt.Errorf("the config path variable %s has not been set; please set", ConfigFileEnvVar)
	}

	appCfg := &AppConfig{}
	_, err := configuration.NewConfiguration(
		configuration.WithType("yaml"),
		configuration.WithName(ConfigurationName),
		configuration.WithReader(cfgFileReader),
		configuration.WithData(appCfg))

	if err != nil {
		return nil, err
	}

	return appCfg, nil
}

func (m *AppConfig) GetDefaults() []configuration.VarDefinition {
	vd := make([]configuration.VarDefinition, 0, 20)
	vd = append(vd, GetHttpSrvConfigDefaults()...)
	vd = append(vd, GetMiddlewareConfigDefaults("config.mw-handler-registry")...)
	return vd
}

func GetHttpSrvConfigDefaults() []configuration.VarDefinition {
	return []configuration.VarDefinition{
		{"config.http.bind-address", httpsrv.DefaultBindAddress, "host reference"},
		{"config.http.server-context.path", httpsrv.DefaultContextPath, "context-path"},
		{"config.http.port", httpsrv.DefaultListenPort, "port"},
		{"config.http.shutdown-timeout", httpsrv.DefaultShutdownTimeout, "shutdown timeout"},
		{"config.http.server-mode", httpsrv.DefaultServerMode, "modalita' di lavoro server gin"},
	}
}

func GetMiddlewareConfigDefaults(contextPath string) []configuration.VarDefinition {
	return []configuration.VarDefinition{
		{strings.Join([]string{contextPath, middleware.MiddlewareErrorId, "disclose-error-info"}, "."), middleware.MiddlewareErrorDefaultDiscoleInfo, "error is in clear"},
		{strings.Join([]string{contextPath, middleware.MiddlewareTracingId, "alphabet"}, "."), middleware.MiddlewareTracingDefaultAlphabet, "alphabet"},
		{strings.Join([]string{contextPath, middleware.MiddlewareTracingId, "spantag"}, "."), middleware.MiddlewareTracingDefaultSpanTag, "spantag"},
		{strings.Join([]string{contextPath, middleware.MiddlewareTracingId, "header"}, "."), middleware.MiddlewareTracingDefaultHeader, "header"},
		{strings.Join([]string{contextPath, middleware.MiddlewareMetricsPromHttpId, "header"}, "."), middleware.MiddlewareTracingDefaultHeader, "header"},
	}
}

func (m *AppConfig) PostProcess() error {

	return nil
}

package main

import (
	_ "embed"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-gin/httpsrv"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-gin/middleware"
	"github.com/rs/zerolog/log"
	"gitlab.alm.poste.it/go/observability/tracing"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tpm-symphony/api"
	"tpm-symphony/registry"
)

//go:embed sha.txt
var sha string

//go:embed VERSION
var version string

// appLogo contains the ASCII splash screen
//go:embed app-logo.txt
var appLogo []byte

func main() {

	fmt.Println(string(appLogo))
	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Sha: %s\n", sha)
	appCfg, err := ReadConfig()
	if nil != err {
		log.Fatal().Err(err).Send()
	}

	log.Info().Interface("config", appCfg).Msg("configuration loaded")

	reg, err := registry.LoadRegistry(&appCfg.Registry)
	if nil != err {
		log.Fatal().Err(err).Send()
	}
	reg.ShowInfo()

	api.RegisterEndpoints(reg)

	jc, err := InitGlobalTracer()
	if nil != err {
		log.Fatal().Err(err).Send()
	}
	defer jc.Close()

	if appCfg.MwRegistry != nil {
		if err := middleware.InitializeHandlerRegistry(appCfg.MwRegistry); err != nil {
			log.Fatal().Err(err).Send()
		}
	}

	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	s, err := httpsrv.NewServer(appCfg.Http, httpsrv.WithListenPort(9090), httpsrv.WithDocumentRoot("/www", "/tmp", false))
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	if err := s.Start(); err != nil {
		log.Fatal().Err(err).Send()
	}
	defer s.Stop()

	for !s.IsReady() {
		time.Sleep(time.Duration(500) * time.Millisecond)
	}

	sig := <-shutdownChannel
	log.Debug().Interface("signal", sig).Msg("got termination signal")

}

func InitGlobalTracer() (*tracing.Tracer, error) {
	tracer, err := tracing.NewTracer()
	if err != nil {
		return nil, err
	}

	return tracer, err
}

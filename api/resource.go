package api

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-gin/httpsrv"
	"github.com/gin-gonic/gin"
	varResolver "github.com/mario-imperato/tpm-common/util/vars"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"tpm-symphony/registry"
)

func RegisterEndpoints(registry *registry.OrchestrationRegistry) {
	log.Info().Msg("registering defined endpoints")
	ra := httpsrv.GetApp()
	ra.RegisterGFactory(registerEndpoints(registry))
}

func registerEndpoints(registry *registry.OrchestrationRegistry) httpsrv.GFactory {
	return func(ctx httpsrv.ServerContext) []httpsrv.G {

		hgs := make([]httpsrv.G, 0, len(registry.ResourceGroups))

		for _, rg := range registry.ResourceGroups {

			gp := rg.Path
			gp = strings.TrimPrefix(gp, ctx.GetContextPath())
			if gp[:1] == "/" {
				gp = gp[1:]
			}

			hg := httpsrv.G{
				Name: rg.Name,
				Path: gp,
			}

			for _, r := range rg.Resources {

				ginVarFormatter := func(s string) string { return ":" + s }
				ginPath, _ := varResolver.ResolveVariables(r.Path, varResolver.SimpleVariableReference, ginVarFormatter)

				hg.Resources = append(hg.Resources, httpsrv.R{
					Name:          r.Name,
					Path:          ginPath,
					Method:        r.Method,
					RouteHandlers: []httpsrv.H{executeOrchestation(r.SymphonyId)},
				})
			}

			hgs = append(hgs, hg)
		}

		return hgs
	}
}

func executeOrchestation(sid string) httpsrv.H {
	return func(c *gin.Context) {
		c.String(http.StatusOK, fmt.Sprintf("Hello Orchestration %s", sid))
	}
}

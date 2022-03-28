package api

import (
	"context"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-gin/httpsrv"
	"github.com/getkin/kin-openapi/openapi3filter"
	legacyrouter "github.com/getkin/kin-openapi/routers/legacy"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/mario-imperato/tpm-common/util"
	varResolver "github.com/mario-imperato/tpm-common/util/vars"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"time"
	"tpm-symphony/registry"
	"tpm-symphony/util/oapiutil"
)

func RegisterEndpoints(registry *registry.OrchestrationRegistry) {
	log.Info().Msg("registering defined endpoints")
	ra := httpsrv.GetApp()
	ra.RegisterGFactory(registerEndpoints(registry))
}

func registerEndpoints(registry *registry.OrchestrationRegistry) httpsrv.GFactory {
	return func(ctx httpsrv.ServerContext) []httpsrv.G {

		hgs := make([]httpsrv.G, 0)
		for _, e := range registry.Entries {

			rg := e.ResourceGroup
			gp := rg.Path
			gp = strings.TrimPrefix(gp, ctx.GetContextPath())
			if len(gp) > 0 && gp[:1] == "/" {
				gp = gp[1:]
			}

			hg := httpsrv.G{
				Name: rg.Name,
				Path: gp,
			}

			for _, r := range e.Resources {

				r.Path = strings.TrimPrefix(r.Path, rg.Path)
				ginVarFormatter := func(s string) string { return ":" + s }
				ginPath, _ := varResolver.ResolveVariables(r.Path, varResolver.SimpleVariableReference, ginVarFormatter)

				hg.Resources = append(hg.Resources, httpsrv.R{
					Name:          r.Name,
					Path:          ginPath,
					Method:        r.Method,
					RouteHandlers: []httpsrv.H{executeOrchestation(&e, r.Path /*path.Join(ctx.GetContextPath(), gp, r.Path) */, r.SymphonyId)},
				})
			}

			hgs = append(hgs, hg)
		}

		return hgs
	}
}

func executeOrchestation(entry *registry.Entry, oapiPath string, sid string) httpsrv.H {
	return func(c *gin.Context) {

		doc := entry.OpenapiDoc
		router, _ := legacyrouter.NewRouter(doc)

		// Find route
		route, pathParams, _ := router.FindRoute(c.Request)

		// Validate request
		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    c.Request,
			PathParams: pathParams,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(context.Background(), requestValidationInput); err != nil {
			produceValidationFailedResponse(c, entry, oapiPath, err)
			return
		}

		c.String(http.StatusOK, fmt.Sprintf("Hello Orchestration %s", sid))
	}
}

func produceValidationFailedResponse(c *gin.Context, entry *registry.Entry, oapiPath string, validationErr error) {

	log.Error().Err(oapiutil.ValidationError(validationErr)).Msg("request validation failed")

	var response []byte
	var err error
	var externalRef string

	response, externalRef, err = oapiutil.RetrieveExample(entry.OpenapiDoc, c.Request.Method, oapiPath, 400, "application/json")
	if err != nil {
		// Should handle built-in response
		response = []byte("{ \"err\": -99 }")
	} else {
		if externalRef != "" {
			response, err = entry.Repo.GetOpenApiAssets(externalRef)
			if err != nil {
				log.Error().Err(err).Str("external-ref", externalRef).Msg("cannot find external reference")
				response = []byte("{ \"err\": -99 }")
			}
		}
	}

	m := map[string]interface{}{
		"message": util.JSONEscape(oapiutil.ValidationError(validationErr).Error()),
		"ts":      time.Now().Format(time.RFC3339),
	}

	t, err := util.ParseTemplates([]util.TemplateInfo{{Name: "validation", Content: string(response)}}, nil)
	if err != nil {
		response = []byte("{ \"err\": -99 }")
	}

	response, err = util.ProcessTemplate(t, m, false)
	if err != nil {
		response = []byte("{ \"err\": -99 }")
	}

	respJson := make(map[string]interface{})
	err = jsoniter.Unmarshal(response, &respJson)
	if err != nil {
		response = []byte("{ \"err\": -99 }")
	}
	c.JSON(http.StatusBadRequest, respJson)
}

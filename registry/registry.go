package registry

import (
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"tpm-symphony/registry/crawler"
)

var theRegistry OrchestrationRegistry

func GetRegistry() *OrchestrationRegistry {
	return &theRegistry
}

func LoadRegistry(cfg *Config) (*OrchestrationRegistry, error) {

	theRegistry = OrchestrationRegistry{}

	oapi, _, err := crawler.Crawl(&cfg.CrawlerCfg)
	if err != nil {
		return nil, err
	}

	if oapi.IsZero() {
		return nil, fmt.Errorf("no orchestrations found in %s", cfg.CrawlerCfg.MountPoint)
	}

	log.Info().Str("file-name", oapi.Path).Msg("found openapi file")

	doc, err := openapi3.NewLoader().LoadFromFile(oapi.Path)
	if err != nil {
		return nil, err
	}

	theRegistry.OpenapiDoc = doc
	theRegistry.OpenApiFile = oapi.Path
	for pi, pv := range doc.Paths {
		pathInfo, oinfos := retrieveSymphonyPathExtension(pi, pv)
		if len(oinfos) == 0 {
			log.Warn().Str("url", pi).Msg("openapi info without symphony information")
		} else {
			for _, op := range oinfos {
				theRegistry.AddResource(
					ResourceGroup{
						Name: pathInfo.GroupPath,
						Path: pathInfo.GroupPath,
					},
					Resource{
						Name:       op.Id,
						Path:       pathInfo.OpPath,
						Method:     op.HttpMethod,
						SymphonyId: op.Id,
					})
			}
		}
	}

	return &theRegistry, nil
}

func retrieveSymphonyPathExtension(pi string, pv *openapi3.PathItem) (SymphonyPathExtension, []SymphonyOperationExtension) {

	var pathExt SymphonyPathExtension
	if pv != nil && pv.Extensions != nil {
		v, ok := pv.Extensions["x-symphony"]
		if ok {
			jr, ok := v.(json.RawMessage)
			if ok {
				err := jsoniter.Unmarshal(jr, &pathExt)
				if err != nil {
					log.Error().Err(err).Msg("wrong sid")
				}
			}
		}
	}

	pathExt.OpPath = pi
	if len(pathExt.GroupPath) > 0 && pathExt.GroupPath != "/" {
		if strings.HasPrefix(pi, pathExt.GroupPath) {
			pathExt.OpPath = strings.TrimPrefix(pi, pathExt.GroupPath)
		}
	} else if strings.HasPrefix(pi, "/api/v1") {
		pathExt.GroupPath = "/api/v1"
		pathExt.OpPath = strings.TrimPrefix(pi, pathExt.GroupPath)
	}

	var oinfo []SymphonyOperationExtension
	oi := retrieveSymphonyOperationExtension(http.MethodGet, pv.Get)
	if !oi.IsZero() {
		oinfo = append(oinfo, oi)
	} else {
		warnForMissingSymphonyInfo(http.MethodGet, pi, pv.Get)
	}
	oi = retrieveSymphonyOperationExtension(http.MethodPut, pv.Put)
	if !oi.IsZero() {
		oinfo = append(oinfo, oi)
	} else {
		warnForMissingSymphonyInfo(http.MethodPut, pi, pv.Put)
	}
	oi = retrieveSymphonyOperationExtension(http.MethodPost, pv.Post)
	if !oi.IsZero() {
		oinfo = append(oinfo, oi)
	} else {
		warnForMissingSymphonyInfo(http.MethodPost, pi, pv.Post)
	}
	oi = retrieveSymphonyOperationExtension(http.MethodDelete, pv.Delete)
	if !oi.IsZero() {
		oinfo = append(oinfo, oi)
	} else {
		warnForMissingSymphonyInfo(http.MethodDelete, pi, pv.Delete)
	}

	return pathExt, oinfo
}

func warnForMissingSymphonyInfo(httpMethod string, path string, op *openapi3.Operation) {
	if op != nil {
		log.Warn().Str("path", path).Str("http-method", httpMethod).Msg("the method has no matching orchestration")
	}
}

func retrieveSymphonyOperationExtension(method string, op *openapi3.Operation) SymphonyOperationExtension {

	if op != nil && op.Extensions != nil {
		v, ok := op.Extensions["x-symphony"]
		if ok {
			jr, ok := v.(json.RawMessage)
			if ok {
				var sid SymphonyOperationExtension
				err := json.Unmarshal(jr, &sid)
				if err == nil {
					sid.HttpMethod = method
					return sid
				}

				log.Error().Err(err).Msg("wrong sid")
			}
		}
	}

	return SymphonyOperationExtension{}
}

func (reg *OrchestrationRegistry) AddResource(group ResourceGroup, res Resource) {

	gndx := reg.FindGroupByPath(group.Path)
	if gndx < 0 {
		reg.ResourceGroups = append(reg.ResourceGroups, group)
		gndx = 0
	}

	reg.ResourceGroups[gndx].Resources = append(reg.ResourceGroups[gndx].Resources, res)
}

func (reg *OrchestrationRegistry) FindGroupByPath(p string) int {
	for i := range reg.ResourceGroups {
		if reg.ResourceGroups[i].Path == p {
			return i
		}
	}

	return -1
}

func (reg *OrchestrationRegistry) ShowInfo() {
	log.Info().Str("from", reg.OpenApiFile).Str("open-api-ver", reg.OpenapiDoc.OpenAPI).Msg("loaded openapi spec")

	for _, g := range reg.ResourceGroups {
		log.Info().Int("num-resources", len(g.Resources)).Str("path", g.Path).Str("name", g.Name).Msg("resource group")
		for _, r := range g.Resources {
			/*
				Bad in log string....
				ginVarFormatter := func(s string) string { return ":" + s }
				ginPath, _ := varResolver.ResolveVariables(r.Path, varResolver.SimpleVariableReference, ginVarFormatter)
			*/
			log.Info().Str("path", r.Path).Str("name", r.Name).Str("method", r.Method).Str("sid", r.SymphonyId).Msg("openapi path references orchestrations")
		}
	}
}

package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/rs/zerolog/log"
	"net/http"
	"regexp"
	"tpm-symphony/constants"
	"tpm-symphony/registry/crawler"
)

var theRegistry OrchestrationRegistry

func GetRegistry() *OrchestrationRegistry {
	return &theRegistry
}

func LoadRegistry(cfg *Config) (*OrchestrationRegistry, error) {

	theRegistry = OrchestrationRegistry{}

	orchs, err := crawler.Crawl(&cfg.CrawlerCfg)
	if err != nil {
		return nil, err
	}

	if len(orchs) == 0 {
		return nil, fmt.Errorf("no orchestrations found in %s", cfg.CrawlerCfg.MountPoint)
	}

	for _, repo := range orchs {

		oapiName, oapiContent, err := repo.GetOpenApiData()
		log.Info().Str(constants.SemLogPath, repo.GetPath()).Str("name", oapiName).Msg("found openapi file")
		doc, err := openapi3.NewLoader().LoadFromData(oapiContent)
		if err != nil {
			return nil, err
		}

		err = validateOpenapiDoc(doc)
		if err != nil {
			return nil, err
		}

		entry := Entry{OpenapiDoc: doc, Repo: repo, ResourceGroup: retrieveResourceGroup(doc)}
		for pi, pv := range doc.Paths {
			oinfos := retrieveSymphonyPathExtension(pi, pv)
			if len(oinfos) == 0 {
				log.Warn().Str("url", pi).Msg("openapi info without symphony information...skipping")
			} else {
				for _, op := range oinfos {

					entry.AddResource(
						Resource{
							Name:       op.Id,
							Path:       pi,
							Method:     op.HttpMethod,
							SymphonyId: op.Id,
						})
				}
			}
		}

		theRegistry.Entries = append(theRegistry.Entries, entry)
	}

	return &theRegistry, nil
}

var ServersUrlPattern = regexp.MustCompile(`^(?:http|https)://[0-9a-zA-Z\.]*(?:\:[0-9]{2,4})?(.*)`)

func validateOpenapiDoc(doc *openapi3.T) error {
	err := doc.Validate(context.Background())
	if err != nil {
		return err
	}

	// Hack. In case the url is in the form http://localhost:8080/something... I reset to /something....
	// This because otherwise the openapi3 lib doesn't match it.
	if len(doc.Servers) > 0 {
		u := util.ExtractCapturedGroupIfMatch(ServersUrlPattern, doc.Servers[0].URL)
		doc.Servers[0].URL = u
	}

	return nil
}

func retrieveResourceGroup(doc *openapi3.T) ResourceGroup {
	if len(doc.Servers) > 0 {
		u := doc.Servers[0].URL
		return ResourceGroup{Name: u, Path: u}
	}

	return ResourceGroup{}
}

func retrieveSymphonyPathExtension(pi string, pv *openapi3.PathItem) []SymphonyOperationExtension {

	/*
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
	*/

	var oinfo []SymphonyOperationExtension
	httpMethods := []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete}
	for _, m := range httpMethods {
		op := pv.GetOperation(m)
		oi := retrieveSymphonyOperationExtension(m, op)
		if !oi.IsZero() {
			oinfo = append(oinfo, oi)
		} else {
			warnForMissingSymphonyInfo(m, pi, op)
		}
	}

	return oinfo
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

func (reg *Entry) AddResource(res Resource) {

	reg.Resources = append(reg.Resources, res)
}

/*
func (reg *Entry) FindGroupByPath(p string) int {
	for i := range reg.ResourceGroups {
		if reg.ResourceGroups[i].Path == p {
			return i
		}
	}

	return -1
}
*/

func (reg *OrchestrationRegistry) ShowInfo() {

	for _, e := range reg.Entries {
		e.Repo.ShowInfo()
	}

	for _, e := range reg.Entries {
		log.Info().Str(constants.SemLogPath, e.Repo.GetPath()).Str("open-api-ver", e.OpenapiDoc.OpenAPI).Msg("open api info")
		g := e.ResourceGroup
		log.Info().Int("num-resources", len(e.Resources)).Str(constants.SemLogPath, g.Path).Str(constants.SemLogName, g.Name).Msg("resource group")
		for _, r := range e.Resources {
			log.Info().Str(constants.SemLogPath, r.Path).Str(constants.SemLogName, r.Name).Str(constants.SemLogMethod, r.Method).Str(constants.SemLogOrchestrationSid, r.SymphonyId).Msg("openapi path references orchestrations")
		}
	}
}

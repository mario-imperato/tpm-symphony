package executable

import (
	"github.com/rs/zerolog/log"
	"tpm-symphony/orchestration/config"
)

type Executable interface {
	Execute() error
	AddInput(p Path) error
	AddOutput(p Path) error
	IsValid() bool
}

type Activity struct {
	Cfg     config.Configurable
	Outputs []Path
	Inputs  []Path
}

func (a *Activity) AddOutput(p Path) error {
	a.Outputs = append(a.Outputs, p)
	return nil
}

func (a *Activity) AddInput(p Path) error {
	a.Inputs = append(a.Inputs, p)
	return nil
}

func (a *Activity) IsValid() bool {

	rc := true
	switch a.Cfg.Type() {
	case config.StartActivityType:
		if len(a.Outputs) == 0 {
			log.Trace().Str("name", a.Cfg.Name()).Int("len-outputs", len(a.Outputs)).Msg("start activity missing outputs")
			rc = false
		}

		if len(a.Inputs) != 0 {
			log.Trace().Str("name", a.Cfg.Name()).Int("len-inputs", len(a.Inputs)).Msg("start activity doesn't have inputs")
			rc = false
		}

	case config.EndActivityType:
		if len(a.Inputs) == 0 {
			log.Trace().Str("name", a.Cfg.Name()).Int("len-inputs", len(a.Inputs)).Msg("end activity missing inputs")
			rc = false
		}

		if len(a.Outputs) != 0 {
			log.Trace().Str("name", a.Cfg.Name()).Int("len-outputs", len(a.Outputs)).Msg("end activity doesn't have outputs")
			rc = false
		}

	default:
		if len(a.Inputs) == 0 || len(a.Outputs) == 0 {
			log.Trace().Str("name", a.Cfg.Name()).Int("len-outputs", len(a.Outputs)).Int("len-inputs", len(a.Inputs)).Msg("activity missing connections")
			rc = false
		}
	}

	return rc
}

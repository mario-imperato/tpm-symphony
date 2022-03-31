package executable

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"tpm-symphony/orchestration/config"
)

type Orchestration struct {
	Cfg         *config.Orchestration
	Executables map[string]Executable
}

func NewOrchestration(cfg *config.Orchestration) (Orchestration, error) {

	o := Orchestration{Cfg: cfg}
	var execs map[string]Executable

	for _, cfgItem := range cfg.Activities {

		var ex Executable
		var err error
		switch cfgItem.Type() {
		case config.StartActivityType:
			ex, err = NewStartActivity(cfgItem)
		case config.EchoActivityType:
			ex, err = NewEchoActivity(cfgItem)
		case config.EndActivityType:
			ex, err = NewEndActivity(cfgItem)
		default:
			panic(fmt.Errorf("this should not happen %s, unrecognized sctivity type", cfgItem.Type()))
		}

		if err != nil {
			return o, err
		}

		if execs == nil {
			execs = make(map[string]Executable)
		}

		execs[cfgItem.Name()] = ex
	}

	if len(execs) == 0 {
		return o, errors.New("empty orchestration found")
	}
	o.Executables = execs

	for _, pcfg := range cfg.Paths {
		p, err := NewPath(&pcfg)
		if err != nil {
			return o, err
		}

		var ex Executable
		var ok bool
		if ex, ok = execs[pcfg.SourceName]; !ok {
			return o, fmt.Errorf("dangling path, could not find source %s", pcfg.SourceName)
		}

		ex.AddOutput(p)

		if ex, ok = execs[pcfg.TargetName]; !ok {
			return o, fmt.Errorf("dangling path, could not find target %s", pcfg.TargetName)
		}

		ex.AddInput(p)
	}

	return o, nil
}

func (o *Orchestration) IsValid() bool {

	if len(o.Executables) == 0 {
		log.Trace().Msg("empty orchestration found")
		return false
	}

	rc := true
	for _, ex := range o.Executables {
		if !ex.IsValid() {
			rc = false
		}
	}

	return rc
}

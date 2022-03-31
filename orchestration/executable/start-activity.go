package executable

import "tpm-symphony/orchestration/config"

type StartActivity struct {
	Activity
}

func NewStartActivity(item config.Configurable) (*StartActivity, error) {

	a := &StartActivity{}
	a.Cfg = item
	return a, nil
}

func (a *StartActivity) Execute() error {
	return nil
}

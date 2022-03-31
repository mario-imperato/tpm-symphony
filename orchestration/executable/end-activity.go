package executable

import "tpm-symphony/orchestration/config"

type EndActivity struct {
	Activity
}

func NewEndActivity(item config.Configurable) (*EndActivity, error) {

	a := &EndActivity{}
	a.Cfg = item
	return a, nil
}

func (a *EndActivity) Execute() error {
	return nil
}

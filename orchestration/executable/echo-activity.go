package executable

import "tpm-symphony/orchestration/config"

type EchoActivity struct {
	Activity
}

func NewEchoActivity(item config.Configurable) (*EchoActivity, error) {

	ea := &EchoActivity{}
	ea.Cfg = item
	return ea, nil
}

func (a *EchoActivity) Execute() error {
	return nil
}

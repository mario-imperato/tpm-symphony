package config

type EchoActivity struct {
	Activity
	Message string `yaml:"message,omitempty" mapstructure:"message,omitempty" json:"message,omitempty"`
}

func (c *EchoActivity) WithName(n string) *EchoActivity {
	c.Nm = n
	return c
}

func NewEchoActivity() *EchoActivity {
	s := EchoActivity{}
	s.Tp = EchoActivityType
	return &s
}

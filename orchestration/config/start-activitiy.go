package config

type StartActivity struct {
	Activity
	StartAProperty string `yaml:"property,omitempty" mapstructure:"property,omitempty" json:"property,omitempty"`
}

func (c *StartActivity) WithName(n string) *StartActivity {
	c.Nm = n
	return c
}

func NewStartActivity() *StartActivity {
	s := StartActivity{}
	s.Tp = StartActivityType
	return &s
}

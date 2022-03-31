package config

type EndActivity struct {
	Activity
	EndAProperty string `yaml:"property,omitempty" mapstructure:"property,omitempty" json:"property,omitempty"`
}

func (c *EndActivity) WithName(n string) *EndActivity {
	c.Nm = n
	return c
}

func NewEndActivity() *EndActivity {
	a := EndActivity{}
	a.Tp = EndActivityType
	return &a
}

package config

type Type string

const (
	StartActivityType Type = "start-activity"
	EchoActivityType  Type = "echo-activity"
	EndActivityType   Type = "end-activity"
)

type Configurable interface {
	Name() string
	Type() Type
}

type Activity struct {
	Nm string `yaml:"name,omitempty" mapstructure:"name,omitempty" json:"name,omitempty"`
	Tp Type   `yaml:"type,omitempty" mapstructure:"type,omitempty" json:"type,omitempty"`
}

func (c *Activity) WithName(n string) *Activity {
	c.Nm = n
	return c
}

func (c *Activity) Name() string {
	return c.Nm
}

func (c *Activity) Type() Type {
	return c.Tp
}

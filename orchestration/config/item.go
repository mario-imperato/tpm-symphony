package config

type ActivityItem interface {
	Name() string
	ActivityType() string
}

type Item struct {
	Nm string `yaml:"name,omitempty" mapstructure:"name,omitempty" json:"name,omitempty"`
	Tp string `yaml:"type,omitempty" mapstructure:"type,omitempty" json:"type,omitempty"`
}

func (c *Item) WithName(n string) *Item {
	c.Nm = n
	return c
}

func (c *Item) Name() string {
	return c.Nm
}

func (c *Item) ActivityType() string {
	return c.Tp
}

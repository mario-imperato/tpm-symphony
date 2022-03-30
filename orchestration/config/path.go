package config

type Path struct {
	SourceName string `yaml:"source,omitempty" mapstructure:"source,omitempty" json:"source,omitempty"`
	TargetName string `yaml:"target,omitempty" mapstructure:"target,omitempty" json:"target,omitempty"`
}

func NewPath(source string, target string) *Path {
	p := Path{SourceName: source, TargetName: target}
	return &p
}

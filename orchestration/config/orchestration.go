package config

import (
	"encoding/json"
	"fmt"
)

type Orchestration struct {
	Paths         []Path            `yaml:"paths,omitempty" mapstructure:"paths,omitempty" json:"paths,omitempty"`
	Activities    []ActivityItem    `json:"-"`
	RawActivities []json.RawMessage `json:"activities"`
}

func (o *Orchestration) ToJSON() ([]byte, error) {
	return json.Marshal(o)
}

func (o *Orchestration) FindActivityByName(n string) ActivityItem {
	for _, a := range o.Activities {
		if a.Name() == n {
			return a
		}
	}

	return nil
}

func (o *Orchestration) AddActivity(a ActivityItem) error {

	if o.FindActivityByName(a.Name()) != nil {
		return fmt.Errorf("activity with the same id already present (id: %s)", a.Name())
	}

	o.Activities = append(o.Activities, a)
	return nil
}

func (o *Orchestration) AddPath(source, target string) error {

	if source == "" || target == "" {
		return fmt.Errorf("path missing source or target reference")
	}

	if o.FindActivityByName(source) == nil {
		return fmt.Errorf("cannot find source activity (id: %s)", source)
	}

	if o.FindActivityByName(target) == nil {
		return fmt.Errorf("cannot find target activity (id: %s)", target)
	}

	o.Paths = append(o.Paths, Path{SourceName: source, TargetName: target})
	return nil
}

func (o *Orchestration) UnmarshalJSON(b []byte) error {

	// Clear the state....
	o.Activities = nil

	type orchestration Orchestration
	err := json.Unmarshal(b, (*orchestration)(o))
	if err != nil {
		return err
	}

	for _, raw := range o.RawActivities {
		var v Item
		err = json.Unmarshal(raw, &v)
		if err != nil {
			return err
		}
		var i ActivityItem
		switch v.ActivityType() {
		case StartActivityType:
			i = NewStartActivity()
		case EchoActivityType:
			i = NewEchoActivity()
		default:
			return fmt.Errorf("unknown activity type %s", v.ActivityType())
		}
		err = json.Unmarshal(raw, i)
		if err != nil {
			return err
		}
		o.Activities = append(o.Activities, i)
	}
	return nil
}

func (o *Orchestration) MarshalJSON() ([]byte, error) {

	// Clear the state....
	o.RawActivities = nil

	type orchestration Orchestration
	if o.Activities != nil {
		for _, v := range o.Activities {
			b, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			o.RawActivities = append(o.RawActivities, b)
		}
	}
	return json.Marshal((*orchestration)(o))
}

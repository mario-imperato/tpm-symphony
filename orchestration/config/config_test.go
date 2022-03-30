package config_test

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
	"tpm-symphony/orchestration/config"
)

func TestConfig(t *testing.T) {

	// Serialization
	sa := config.NewStartActivity().WithName("start-name")
	sa.StartAProperty = "a-start-property"

	ea := config.NewEchoActivity().WithName("echo-name")
	ea.Message = "a-message"

	orch := config.Orchestration{
		Activities: []config.ActivityItem{
			sa, ea,
		},
	}

	err := orch.AddPath("start-name", "echo-name")
	require.NoError(t, err)

	b, err := json.Marshal(&orch)
	require.NoError(t, err)
	t.Log(string(b))

	// Deserialization
	var orch2 = config.Orchestration{}
	err = json.Unmarshal(b, &orch2)
	require.NoError(t, err)

	b, err = orch2.ToJSON()
	require.NoError(t, err)
	t.Log(string(b))
}

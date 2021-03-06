package config_test

import (
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

	ea2 := config.NewEndActivity().WithName("end-name")

	orch := config.Orchestration{
		Activities: []config.Configurable{
			sa, ea, ea2,
		},
	}

	err := orch.AddPath("start-name", "echo-name")
	require.NoError(t, err)

	err = orch.AddPath("echo-name", "end-name")
	require.NoError(t, err)

	b, err := orch.ToJSON()
	require.NoError(t, err)
	t.Log(string(b))

	// Deserialization
	orch2, err := config.NewOrchestration("", b)
	require.NoError(t, err)

	b, err = orch2.ToJSON()
	require.NoError(t, err)
	t.Log(string(b))
}

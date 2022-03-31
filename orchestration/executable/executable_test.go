package executable_test

import (
	"github.com/stretchr/testify/require"
	"testing"
	"tpm-symphony/orchestration/config"
	"tpm-symphony/orchestration/executable"
)

func TestNewOrchestration(t *testing.T) {

	// Serialization
	sa := config.NewStartActivity().WithName("start-name")
	sa.StartAProperty = "a-start-property"

	ea := config.NewEchoActivity().WithName("echo-name")
	ea.Message = "a-message"

	ea2 := config.NewEndActivity().WithName("end-name")

	cfgOrc := config.Orchestration{
		Activities: []config.Configurable{
			sa, ea, ea2,
		},
	}

	err := cfgOrc.AddPath("start-name", "echo-name")
	require.NoError(t, err)

	err = cfgOrc.AddPath("echo-name", "end-name")
	require.NoError(t, err)

	orc, err := executable.NewOrchestration(&cfgOrc)
	require.NoError(t, err)

	if !orc.IsValid() {
		t.Error("orchestration is invalid")
	}
	t.Log(orc)
}

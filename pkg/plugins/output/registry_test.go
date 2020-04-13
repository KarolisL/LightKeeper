package output

import (
	"fmt"
	"github.com/KarolisL/lightkeeper/pkg/common"
	"github.com/KarolisL/lightkeeper/pkg/daemon/config"
	"github.com/KarolisL/lightkeeper/pkg/test_utils"
	"github.com/google/go-cmp/cmp"
	"testing"
)

type call struct {
	Params config.Params
}

type stubOutput struct {
	Name string
}

func (s *stubOutput) Ch() chan<- common.Message {
	panic("implement me")
}

func TestPluginRegistry_NewOutput(t *testing.T) {
	t.Run("When plugin does exist", func(t *testing.T) {
		output := &stubOutput{"myoutput"}
		var calls []call
		stub := func(params config.Params) (Output, error) {
			calls = append(calls, call{params})
			return output, nil
		}
		registry := PluginRegistry{}
		registry.Register("syslog-ng", stub)

		got, _ := registry.NewOutput("syslog-ng", make(config.Params))
		want := output

		assertOutputReturned(t, want, got)
		assertStubCalls(t, calls, []call{{map[string]string{}}})
	})

	t.Run("When plugin params are not valid propagates error from plugin", func(t *testing.T) {
		stubError := fmt.Errorf("StubError")
		var calls []call
		stub := func(params config.Params) (Output, error) {
			calls = append(calls, call{params})
			return nil, stubError
		}
		registry := PluginRegistry{}
		registry.Register("syslog-ng", stub)

		_, err := registry.NewOutput("syslog-ng", make(config.Params))

		test_utils.AssertErrorIs(t, err, stubError)
		assertStubCalls(t, calls, []call{{map[string]string{}}})
	})

	t.Run("When plugin doesn't exist (empty registry)", func(t *testing.T) {
		registry := PluginRegistry{}

		_, err := registry.NewOutput("syslog-ng", make(config.Params))
		test_utils.AssertErrorIs(t, err, ErrPluginTypeNotFound)
	})
}

func assertOutputReturned(t *testing.T, want *stubOutput, got Output) {
	t.Helper()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("NewOutput returned wrong output (-want, +got):\n%s", diff)
	}
}

func assertStubCalls(t *testing.T, got, want []call) {
	t.Helper()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("NewOutput called stubOutput wrong amount of times or with wrong arguments (-want, +got):\n%s", diff)
	}
}

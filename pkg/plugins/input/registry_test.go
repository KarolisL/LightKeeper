package input

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

type stubInput struct {
	Name string
}

func (s *stubInput) Ch() <-chan common.Message {
	panic("implement me")
}

func TestPluginRegistry_NewInput(t *testing.T) {
	t.Run("When plugin does exist", func(t *testing.T) {
		input := &stubInput{"myinput"}
		var calls []call
		stub := func(params config.Params) (Input, error) {
			calls = append(calls, call{params})
			return input, nil
		}
		registry := PluginRegistry{}
		registry.Register("syslog-ng", stub)

		got, _ := registry.NewInput("syslog-ng", make(config.Params))
		want := input

		assertInputReturned(t, want, got)
		assertStubCalls(t, calls, []call{{map[string]string{}}})
	})

	t.Run("When plugin params are not valid propagates error from plugin", func(t *testing.T) {
		stubError := fmt.Errorf("StubError")
		var calls []call
		stub := func(params config.Params) (Input, error) {
			calls = append(calls, call{params})
			return nil, stubError
		}
		registry := PluginRegistry{}
		registry.Register("syslog-ng", stub)

		_, err := registry.NewInput("syslog-ng", make(config.Params))

		test_utils.AssertErrorIs(t, err, stubError)
		assertStubCalls(t, calls, []call{{map[string]string{}}})
	})

	t.Run("When plugin doesn't exist (empty registry)", func(t *testing.T) {
		registry := PluginRegistry{}

		_, err := registry.NewInput("syslog-ng", make(config.Params))
		test_utils.AssertErrorIs(t, err, ErrPluginTypeNotFound)
	})
}

func assertInputReturned(t *testing.T, want *stubInput, got Input) {
	t.Helper()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("NewInput returned wrong input (-want, +got):\n%s", diff)
	}
}

func assertStubCalls(t *testing.T, got, want []call) {
	t.Helper()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("NewInput called stubInput wrong amount of times or with wrong arguments (-want, +got):\n%s", diff)
	}
}

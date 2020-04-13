package daemon

import (
	"errors"
	"github.com/KarolisL/lightkeeper/pkg/common"
	"github.com/KarolisL/lightkeeper/pkg/daemon/config"
	"github.com/KarolisL/lightkeeper/pkg/plugins/input"
	"github.com/KarolisL/lightkeeper/pkg/plugins/output"
	"github.com/KarolisL/lightkeeper/pkg/test_utils"
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

type call struct {
	Type   string
	Params map[string]string
}

type stubInputPluginRegistry struct {
	returnValue struct {
		i   input.Input
		err error
	}
	calls []call
}

func (ipr *stubInputPluginRegistry) NewInput(inputType string, params config.Params) (input.Input, error) {
	ipr.calls = append(ipr.calls, call{inputType, params})
	return ipr.returnValue.i, ipr.returnValue.err
}

type stubOutputPluginRegistry struct {
	returnValue struct {
		o   output.Output
		err error
	}
	calls []call
}

func (opr *stubOutputPluginRegistry) NewOutput(outputType string, params config.Params) (output.Output, error) {
	opr.calls = append(opr.calls, call{outputType, params})
	return opr.returnValue.o, opr.returnValue.err
}

func TestNewDaemon(t *testing.T) {
	t.Run("Single input and Output", func(t *testing.T) {
		cfg := makeSimpleConfig()
		ipr := &stubInputPluginRegistry{}
		opr := &stubOutputPluginRegistry{}
		NewDaemon(cfg, ipr, opr)

		assertInputRegistryCalled(t, ipr, []call{
			{
				"syslog-ng",
				map[string]string{"path": "/var/log/messages"}},
		})
		assertOutputRegistryCalled(t, opr, []call{
			{
				"telegram",
				config.Params{"token": "TheToken", "chatId": "3344"},
			}})
	})

	t.Run("Two inputs and two outputs", func(t *testing.T) {
		cfg := &config.Config{
			Inputs: map[string]config.Input{
				"1": makeInput("someInput1", nil),
				"2": makeInput("someInput2", nil),
			},
			Outputs: map[string]config.Output{
				"1": makeOutput("someOutput1", nil),
				"2": makeOutput("someOutput2", nil),
				"3": makeOutput("someOutput3", nil),
			},
		}
		ipr := &stubInputPluginRegistry{}
		opr := &stubOutputPluginRegistry{}
		NewDaemon(cfg, ipr, opr)

		assertInputRegistryCalled(t, ipr, []call{
			{"someInput1", nil},
			{"someInput2", nil},
		})
		assertOutputRegistryCalled(t, opr, []call{
			{"someOutput1", nil},
			{"someOutput2", nil},
			{"someOutput3", nil},
		})
	})

	t.Run("Input plugin returns error", func(t *testing.T) {
		cfg := makeSimpleConfig()
		inputCreationError := errors.New("inputCreationError")

		ipr := &stubInputPluginRegistry{returnValue: struct {
			i   input.Input
			err error
		}{i: nil, err: inputCreationError}}

		opr := &stubOutputPluginRegistry{}
		_, err := NewDaemon(cfg, ipr, opr)

		test_utils.AssertErrorIs(t, err, inputCreationError)
		assertInputRegistryCalled(t, ipr, []call{
			{
				"syslog-ng",
				map[string]string{"path": "/var/log/messages"}},
		})
	})
}

type stubInputPlugin struct {
	ch <-chan common.Message
}

func (s *stubInputPlugin) Ch() <-chan common.Message {
	return s.ch
}

type stubOutputPlugin struct {
	ch chan<- common.Message
}

func (s *stubOutputPlugin) Ch() chan<- common.Message {
	return s.ch
}

func TestDaemon_Start(t *testing.T) {
	t.Run("Test connection from intput to output", func(t *testing.T) {
		cfg := &config.Config{
			Inputs: map[string]config.Input{
				"1": makeInput("someInput1", nil),
				"2": makeInput("someInput2", nil),
			},
			Mappings: []config.Mapping{
				{
					From:    "1",
					To:      "2",
					Filters: nil,
				},
			},
			Outputs: map[string]config.Output{
				"1": makeOutput("someOutput1", nil),
				"2": makeOutput("someOutput2", nil),
				"3": makeOutput("someOutput3", nil),
			},
		}
		inCh := make(chan common.Message)
		outCh := make(chan common.Message)
		stubInput := &stubInputPlugin{inCh}
		stubOutput := &stubOutputPlugin{outCh}
		ipr := &stubInputPluginRegistry{returnValue: struct {
			i   input.Input
			err error
		}{i: stubInput, err: nil}}
		opr := &stubOutputPluginRegistry{returnValue: struct {
			o   output.Output
			err error
		}{o: stubOutput, err: nil}}

		daemon, err := NewDaemon(cfg, ipr, opr)
		if err != nil {
			t.Fatalf("NewDaemon returned error, got %q", err)
		}

		go daemon.Start()
		sendWithTimeout(t, inCh, common.Message("Hi!"))

		got := test_utils.ReceiveWithTimeout(t, outCh)
		want := common.Message("Hi!")

		if got != want {
			t.Errorf("Daemon.Start cause wrong message to be sent, got %q, want %q", got, want)
		}
	})
}

func sendWithTimeout(t *testing.T, ch chan common.Message, message common.Message) {
	t.Helper()
	timeout := 100 * time.Millisecond
	select {
	case <-time.After(timeout):
		t.Fatalf("Wasn't able to send message in %d ms", timeout.Milliseconds())
	case ch <- message:
	}
}

func assertOutputRegistryCalled(t *testing.T, opr *stubOutputPluginRegistry, calls []call) {
	t.Helper()
	if len(opr.calls) != len(calls) {
		t.Errorf("OutputPluginRegistry called wrong amount of times: got %d want %d", len(opr.calls), len(calls))
	}

	for i, want := range calls {
		if diff := cmp.Diff(want, opr.calls[i]); diff != "" {
			t.Errorf("Call to OutputPluginRegistry #%d mismatch (-want +got):\n%s", i, diff)
		}
	}
}

func assertInputRegistryCalled(t *testing.T, ipr *stubInputPluginRegistry, calls []call) {
	t.Helper()
	if len(ipr.calls) != len(calls) {
		t.Errorf("InputPluginRegistry called wrong amount of times: got %d want %d", len(ipr.calls), len(calls))
	}

	for i, want := range calls {
		if diff := cmp.Diff(want, ipr.calls[i]); diff != "" {
			t.Errorf("Call to InputPluginRegistry #%d mismatch (-want +got):\n%s", i, diff)
		}
	}
}

func makeSimpleConfig() *config.Config {
	return &config.Config{
		Inputs: map[string]config.Input{
			"varlogmessages": makeInput("syslog-ng", config.Params{"path": "/var/log/messages"}),
		},
		Mappings: []config.Mapping{
			{
				From: "varlogmessages",
				To:   "telegram",
				Filters: []config.Params{
					{"program": "sshd"},
				},
			},
		},
		Outputs: map[string]config.Output{
			"telegram": makeOutput("telegram", config.Params{
				"token":  "TheToken",
				"chatId": "3344",
			}),
		},
	}
}

func makeInput(typ string, params config.Params) config.Input {
	return config.Input{Type: typ, Params: params}
}

func makeOutput(typ string, params config.Params) config.Output {
	return config.Output{Type: typ, Params: params}
}

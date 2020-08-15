package daemon

import (
	"errors"
	"github.com/KarolisL/lightkeeper/pkg/common"
	"github.com/KarolisL/lightkeeper/pkg/daemon/config"
	"github.com/KarolisL/lightkeeper/pkg/plugins/input"
	mock_input "github.com/KarolisL/lightkeeper/pkg/plugins/input/mock"
	mock_output "github.com/KarolisL/lightkeeper/pkg/plugins/output/mock"
	"github.com/KarolisL/lightkeeper/pkg/test_utils"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestNewDaemon(t *testing.T) {
	t.Run("Single input and Output", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ipr := mock_input.NewMockMaker(ctrl)
		opr := mock_output.NewMockMaker(ctrl)

		ipr.EXPECT().NewFanOutInput(gomock.Eq("syslog-ng"), gomock.Eq(config.Params{"path": "/var/log/messages"}))
		opr.EXPECT().NewOutput(gomock.Eq("telegram"), gomock.Eq(config.Params{"token": "TheToken", "chatId": "3344"}))

		cfg := makeSimpleConfig()
		_, err := NewDaemon(cfg, ipr, opr)
		assertNoError(t, err)
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
		ctrl := gomock.NewController(t)
		ipr := mock_input.NewMockMaker(ctrl)
		opr := mock_output.NewMockMaker(ctrl)

		ipr.EXPECT().NewFanOutInput("someInput1", nil)
		ipr.EXPECT().NewFanOutInput("someInput2", nil)

		opr.EXPECT().NewOutput("someOutput1", nil)
		opr.EXPECT().NewOutput("someOutput2", nil)
		opr.EXPECT().NewOutput("someOutput3", nil)

		_, err := NewDaemon(cfg, ipr, opr)
		assertNoError(t, err)
	})

	t.Run("Input plugin returns error", func(t *testing.T) {
		cfg := makeSimpleConfig()
		inputCreationError := errors.New("inputCreationError")

		ctrl := gomock.NewController(t)
		ipr := mock_input.NewMockMaker(ctrl)
		opr := mock_output.NewMockMaker(ctrl)

		ipr.EXPECT().NewFanOutInput("syslog-ng", config.Params{"path": "/var/log/messages"}).
			Return(nil, inputCreationError)

		_, err := NewDaemon(cfg, ipr, opr)
		test_utils.AssertErrorIs(t, err, inputCreationError)
	})
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Expected NewDaemon to not return error, got %v", err)
	}
}

func TestDaemon_Start(t *testing.T) {
	t.Run("Test connection from input to output", func(t *testing.T) {
		cfg, ctrl, ipr, opr := setupForDaemon(t, config.Mapping{
			From:    "1",
			To:      "2",
			Filters: nil,
		})
		defer ctrl.Finish()

		in1 := &stubInput{ch: make(chan common.Message)}
		ipr.EXPECT().NewFanOutInput(gomock.Eq("i1"), gomock.Any()).Return(input.FanOut(in1), nil)
		in2 := &stubInput{ch: make(chan common.Message)}
		ipr.EXPECT().NewFanOutInput(gomock.Eq("i2"), gomock.Any()).Return(input.FanOut(in2), nil)

		out1 := &stubOutput{ch: make(chan common.Message)}
		opr.EXPECT().NewOutput(gomock.Eq("o1"), gomock.Any()).Return(out1, nil)
		out2 := &stubOutput{ch: make(chan common.Message)}
		opr.EXPECT().NewOutput(gomock.Eq("o2"), gomock.Any()).Return(out2, nil)

		daemon, err := NewDaemon(cfg, ipr, opr)
		assertNoError(t, err)

		go daemon.Start()
		test_utils.SendWithTimeout(t, in1.ch, "Hi!")

		got := test_utils.ReceiveWithTimeout(t, out2.ch)
		want := common.Message("Hi!")

		if got != want {
			t.Errorf("Daemon.Start cause wrong message to be sent, got %q, want %q", got, want)
		}
	})

	t.Run("Test test filtering", func(t *testing.T) {
		cfg, ctrl, ipr, opr := setupForDaemon(t, config.Mapping{
			From: "1",
			To:   "2",
			Filters: []config.Params{{
				"type":    "syslog-ng",
				"program": "sshd",
			}},
		})
		defer ctrl.Finish()

		in1 := &stubInput{ch: make(chan common.Message)}
		ipr.EXPECT().NewFanOutInput(gomock.Eq("i1"), gomock.Any()).Return(input.FanOut(in1), nil)
		in2 := &stubInput{ch: make(chan common.Message)}
		ipr.EXPECT().NewFanOutInput(gomock.Eq("i2"), gomock.Any()).Return(input.FanOut(in2), nil)

		out1 := &stubOutput{ch: make(chan common.Message)}
		opr.EXPECT().NewOutput(gomock.Eq("o1"), gomock.Any()).Return(out1, nil)
		out2 := &stubOutput{ch: make(chan common.Message)}
		opr.EXPECT().NewOutput(gomock.Eq("o2"), gomock.Any()).Return(out2, nil)

		daemon, err := NewDaemon(cfg, ipr, opr)
		assertNoError(t, err)

		go daemon.Start()
		test_utils.SendWithTimeout(t, in1.ch, "Hi!")
		sshdMessage := "Apr 14 09:21:52 some-host sshd[32252]: Accepted publickey for root from 10.10.1.2 port 50919 ssh2: RSA SHA256:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
		test_utils.SendWithTimeout(t, in1.ch, common.Message(sshdMessage))

		got := test_utils.ReceiveWithTimeout(t, out2.ch)
		want := common.Message(sshdMessage)

		if got != want {
			t.Errorf("Daemon.Start cause wrong message to be sent, got %q, want %q", got, want)
		}
	})
}

type stubInput struct {
	ch chan common.Message
}

func (s *stubInput) Ch() <-chan common.Message {
	return s.ch
}

type stubOutput struct {
	ch chan common.Message
}

func (s *stubOutput) Ch() chan<- common.Message {
	return s.ch
}

func TestDemon_Mapping(t *testing.T) {
	cfg, ctrl, ipr, opr := setupForDaemon(
		t,
		config.Mapping{
			From: "1",
			To:   "1",
		}, config.Mapping{
			From: "1",
			To:   "2",
		})

	defer ctrl.Finish()

	in1 := &stubInput{ch: make(chan common.Message)}
	ipr.EXPECT().NewFanOutInput(gomock.Eq("i1"), gomock.Any()).Return(input.FanOut(in1), nil)
	in2 := &stubInput{ch: make(chan common.Message)}
	ipr.EXPECT().NewFanOutInput(gomock.Eq("i2"), gomock.Any()).Return(input.FanOut(in2), nil)

	out1 := &stubOutput{ch: make(chan common.Message)}
	opr.EXPECT().NewOutput(gomock.Eq("o1"), gomock.Any()).Return(out1, nil)
	out2 := &stubOutput{ch: make(chan common.Message)}
	opr.EXPECT().NewOutput(gomock.Eq("o2"), gomock.Any()).Return(out2, nil)

	daemon, err := NewDaemon(cfg, ipr, opr)
	assertNoError(t, err)

	go daemon.Start()

	msg := "Hi!"
	test_utils.SendWithTimeout(t, in1.ch, common.Message(msg))

	got := test_utils.ReceiveWithTimeout(t, out1.ch)
	want := common.Message(msg)

	if got != want {
		t.Fatalf("Wrong message received: got %q, want %q", got, want)
	}

	got = test_utils.ReceiveWithTimeout(t, out2.ch)

	if got != want {
		t.Fatalf("Wrong message received: got %q, want %q", got, want)
	}
}

func setupForDaemon(t *testing.T, mappings ...config.Mapping) (*config.Config, *gomock.Controller, *mock_input.MockMaker, *mock_output.MockMaker) {
	cfg := &config.Config{
		Inputs: map[string]config.Input{
			"1": makeInput("i1", nil),
			"2": makeInput("i2", nil),
		},
		Mappings: mappings,
		Outputs: map[string]config.Output{
			"1": makeOutput("o1", nil),
			"2": makeOutput("o2", nil),
		},
	}

	ctrl := gomock.NewController(t)
	ipr := mock_input.NewMockMaker(ctrl)
	opr := mock_output.NewMockMaker(ctrl)

	return cfg, ctrl, ipr, opr
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

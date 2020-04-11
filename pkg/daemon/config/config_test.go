package config

import (
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
)

// language = toml
const validToml = `
[inputs.varlogmessages]
type = "syslog-ng"
params.path = "/var/log/messages"

[[mappings]]
from = "varlogmessages"
to = "telegram"
filters = [{ program = "sshd" }]

[outputs.telegram]
type = "telegram"
params.token = "TheToken"
params.chatId = "3344"
`

func TestNewConfigFromReader(t *testing.T) {
	t.Run("Parse valid TOML", func(t *testing.T) {

		got, err := NewConfigFromReader(strings.NewReader(validToml))
		want := Config{
			Inputs: map[string]Input{
				"varlogmessages": {
					Type: "syslog-ng",
					Params: map[string]string{
						"path": "/var/log/messages",
					},
				},
			},
			Mappings: []Mapping{
				{
					From: "varlogmessages",
					To:   "telegram",
					Filters: []Params{
						{"program": "sshd"},
					},
				},
			},
			Outputs: map[string]Output{
				"telegram": {
					Type: "telegram",
					Params: map[string]string{
						"token":  "TheToken",
						"chatId": "3344",
					},
				},
			},
		}

		assertNoError(t, err)

		if diff := cmp.Diff(want, *got); diff != "" {
			t.Errorf("NewConfigFromReader mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Returns error on invalid TOML", func(t *testing.T) {
		_, err := NewConfigFromReader(strings.NewReader("invalidToml"))

		assertError(t, err)
	})
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Wanted no error, got %v", err)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("Wanted error, got %q", err)
	}
}

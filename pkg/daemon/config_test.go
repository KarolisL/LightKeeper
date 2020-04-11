package daemon

import (
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
)

// language = toml
const validToml = `
[inputs.varlogmessages]
path = "/var/log/messages"
output = "telegram"
type = "syslog-ng"
filters = [{ program = "sshd" }]

[outputs.telegram]
type = "telegram"
values.token = "TheToken"
values.chatId = "3344"

`

func TestNewConfigFromReader(t *testing.T) {
	t.Run("Parse valid TOML", func(t *testing.T) {

		got, err := NewConfigFromReader(strings.NewReader(validToml))
		want := Config{
			Inputs: map[string]Input{
				"varlogmessages": {
					Path:    "/var/log/messages",
					Type:    "syslog-ng",
					Filters: []Filter{{Program: "sshd"}},
					Output:  "telegram",
				},
			},
			Outputs: map[string]Output{
				"telegram": {
					Type:   "telegram",
					Values: map[string]string{
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

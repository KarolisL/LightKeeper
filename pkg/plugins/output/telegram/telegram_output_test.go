package telegram

import (
	"bytes"
	"fmt"
	"github.com/KarolisL/lightkeeper/pkg/common"
	"github.com/KarolisL/lightkeeper/pkg/test_utils"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestOutput_Ch(t *testing.T) {
	tests := []struct {
		telegramToken  string
		telegramChatId string
		messages       []string
	}{
		{"TheToken", "chat1", []string{"Something have happened"}},
		{"OtherToken", "chat2", []string{"msg1", "msg2"}},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("Sending %v to %q with token %q", test.messages, test.telegramChatId, test.telegramToken), func(t *testing.T) {
			server, bodies, uris := newServer(`OK`)
			defer server.Close()

			telegramToken := test.telegramToken
			chatId := test.telegramChatId

			output, err := NewTelegramOutput(map[string]string{
				"token":    telegramToken,
				"chatId":   chatId,
				"endpoint": server.URL,
			})
			if err != nil {
				t.Fatalf("NewTelegramOutput wanted no error, got %v", err)
			}

			ch := output.Ch()
			for _, message := range test.messages {
				test_utils.SendWithTimeout(t, ch, common.Message(message))
			}

			for _, message := range test.messages {
				uri := getWithTimeout(t, uris)
				assertUri(t, uri, fmt.Sprintf("/bot%s/sendMessage", test.telegramToken))

				body := getWithTimeout(t, bodies)
				assertValuesInBody(t, body, map[string][]string{
					"chat_id": {chatId},
					"text":    {message},
				})
			}
		})
	}
}

func newServer(returnBody string) (*httptest.Server, chan string, chan string) {
	bodies := make(chan string, 1)
	uris := make(chan string, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := bodyAsString(r)
		uris <- r.RequestURI
		bodies <- body
		w.Write([]byte(returnBody))
	}))
	return server, bodies, uris
}

func assertUri(t *testing.T, got string, want string) {
	t.Helper()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Telegram API was called with wrong URI (-want, +got):\n%s", diff)
	}
}

func bodyAsString(r *http.Request) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()
	return body
}

func assertValuesInBody(t *testing.T, body string, want map[string][]string) {
	t.Helper()
	got, err := url.ParseQuery(body)
	if err != nil {
		t.Fatalf("Unable to parse Telegram request body as form params, got %q, err: %v", body, err)
	}
	if diff := cmp.Diff(url.Values(want), got); diff != "" {
		t.Errorf("Telegram API received wrong values (-want, +got):\n%s", diff)
	}
}

func getWithTimeout(t *testing.T, ch <-chan string) (body string) {
	t.Helper()
	timeout := 100 * time.Millisecond
	select {
	case <-time.After(timeout):
		t.Fatalf("Wasn't able to receive message in %d ms", timeout.Milliseconds())
	case body = <-ch:
	}

	return
}

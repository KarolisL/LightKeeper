package test_utils

import (
	"github.com/KarolisL/lightkeeper/pkg/common"
	"testing"
	"time"
)

func ReceiveWithTimeout(t *testing.T, ch <-chan common.Message) common.Message {
	t.Helper()
	timeout := 100 * time.Millisecond
	select {
	case <-time.After(timeout):
		t.Fatalf("Wasn't able to receive message in %d ms", timeout.Milliseconds())
	case message := <-ch:
		return message
	}

	return ""
}

func SendWithTimeout(t *testing.T, ch chan<- common.Message, message common.Message) {
	t.Helper()
	timeout := 100 * time.Millisecond
	select {
	case <-time.After(timeout):
		t.Fatalf("Wasn't able to send message in %d ms", timeout.Milliseconds())
	case ch <- message:
	}
}

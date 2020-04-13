package file

import (
	"fmt"
	"github.com/KarolisL/lightkeeper/pkg/common"
	"github.com/KarolisL/lightkeeper/pkg/test_utils"
	"io/ioutil"
	"os"
	"testing"
)

func TestFileInput_Ch(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{"Pepper", "Pepper"},
		{"Slippy", "Slippy"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file := mkTempFile(t)
			defer os.Remove(file.Name())

			input, err := NewFileInput(map[string]string{
				"path": file.Name(),
			})

			assertNoError(t, err)

			ch := input.Ch()
			fmt.Fprintln(file, test.message)

			got := test_utils.ReceiveWithTimeout(t, ch)
			want := common.Message(test.message)

			assertReceivedMessage(t, got, want)
		})
	}
}

func assertReceivedMessage(t *testing.T, got common.Message, want common.Message) {
	if got != want {
		t.Errorf("Wrong message received from Ch() channel: got %q, want %q", got, want)
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Wanted no error, but got %v", err)
	}
}

func mkTempFile(t *testing.T) *os.File {
	file, err := ioutil.TempFile("", "TestFileInput_Ch")
	if err != nil {
		t.Fatalf("Cannot create temp file")
	}
	return file
}

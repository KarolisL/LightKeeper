package daemon

//
//import (
//	"fmt"
//	"github.com/KarolisL/lightkeeper/pkg/daemon/config"
//	"github.com/KarolisL/lightkeeper/pkg/plugins/input"
//	"io/ioutil"
//	"testing"
//	"time"
//)
//
//func TestDaemonIntegrationTest(t *testing.T) {
//	file, err := ioutil.TempFile("", "TestDaemonIntegrationTest")
//	if err != nil {
//		t.Fatalf("Unable to create temp file")
//	}
//
//	cfg := config.Config{
//		Inputs: map[string]config.Input{
//			"tempfile": makeInput("syslog-ng", config.Params{
//				"path": file.Name(),
//			})},
//		Outputs: map[string]config.Output{
//			"telegram": makeOutput("telegram", config.Params{
//				"token":  "TheToken",
//				"chatId": "3344",
//			}),
//		}}
//
//	inputRegistry := &input.PluginRegistry{}
//	inputRegistry.Register("syslog-ng", nil)
//
//	daemon, err := NewDaemon(cfg, inputRegistry, outputRegistry)
//	if err != nil {
//		t.Fatalf("Unable to create daemon: %+v", err)
//	}
//
//	// Check whether file is beaing read from beginning
//	fmt.Fprintf(file, "Apr 13 13:59:29 router sshd[11824]: Accepted publickey for root from 192.168.88.33 port 61577 ssh2: RSA SHA256:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
//	go daemon.Start()
//	<-time.After(100 * time.Millisecond)
//	assertOutputMockWasCalledWith(t, "Apr 13 13:59:29 router sshd[11824]: Accepted publickey for root from 192.168.88.33 port 61577 ssh2: RSA SHA256:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
//
//	fmt.Fprintf(file, "Apr 14 18:59:29 router sshd[11824]: Accepted publickey for root from 192.168.33.88 port 61577 ssh2: RSA SHA256:BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
//	<-time.After(100 * time.Millisecond)
//	assertOutputMockWasCalledWith(t, "Apr 14 18:59:29 router sshd[11824]: Accepted publickey for root from 192.168.33.88 port 61577 ssh2: RSA SHA256:BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
//}
//

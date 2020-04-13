package stdout

import (
	"github.com/KarolisL/lightkeeper/pkg/common"
	"github.com/KarolisL/lightkeeper/pkg/daemon/config"
	"github.com/KarolisL/lightkeeper/pkg/plugins/output"
)

type Output struct {
	ch chan common.Message
}

func (o *Output) Ch() chan<- common.Message {
	if o.ch == nil {
		o.ch = make(chan common.Message)
		go o.printSync()
	}

	return o.ch
}

func (o *Output) printSync() {
	for msg := range o.ch {
		println(msg)
	}
}

func init() {
	output.Registry.Register("stdout", NewStdoutOutput)
}

func NewStdoutOutput(_ config.Params) (output.Output, error) {
	outp := &Output{}

	return outp, nil
}

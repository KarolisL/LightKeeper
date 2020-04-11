package daemon

import (
	"github.com/KarolisL/lightkeeper/pkg/daemon/config"
	"github.com/KarolisL/lightkeeper/pkg/plugins/input"
	"github.com/KarolisL/lightkeeper/pkg/plugins/output"
)

type Daemon struct {
}

func NewDaemon(config config.Config, inputMaker input.InputMaker, outputMaker output.OutputMaker) *Daemon {
	for _, input := range config.Inputs {
		inputMaker.NewInput(input.Type, input.Params)
	}

	for _, output := range config.Outputs {
		outputMaker.NewOutput(output.Type, output.Params)
	}
	d := &Daemon{}
	return d
}

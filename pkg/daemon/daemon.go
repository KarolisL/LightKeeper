package daemon

import (
	"github.com/KarolisL/lightkeeper/pkg/daemon/config"
	"github.com/KarolisL/lightkeeper/pkg/plugins/input"
	"github.com/KarolisL/lightkeeper/pkg/plugins/output"
)

type Daemon struct {
}

func NewDaemon(config config.Config, inputMaker input.Maker, outputMaker output.OutputMaker) (*Daemon, error) {
	for _, inp := range config.Inputs {
		_, err := inputMaker.NewInput(inp.Type, inp.Params)
		if err != nil {
			return nil, err
		}
	}

	for _, o := range config.Outputs {
		outputMaker.NewOutput(o.Type, o.Params)
	}
	d := &Daemon{}
	return d, nil
}

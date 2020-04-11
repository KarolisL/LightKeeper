package daemon

import "github.com/KarolisL/lightkeeper/pkg/daemon/config"

type Message string

type Input interface {
	Ch() <-chan Message
}

type Output interface {
	Ch() chan<- Message
}

type InputMaker interface {
	NewInput(inputType string, params config.Params) Input
}

type OutputMaker interface {
	NewOutput(outputType string, params config.Params) Output
}

type Daemon struct {
}

func NewDaemon(config config.Config, inputMaker InputMaker, outputMaker OutputMaker) *Daemon {
	for _, input := range config.Inputs {
		inputMaker.NewInput(input.Type, input.Params)
	}

	for _, output := range config.Outputs {
		outputMaker.NewOutput(output.Type, output.Params)
	}
	d := &Daemon{}
	return d
}

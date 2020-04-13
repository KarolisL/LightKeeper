package daemon

import (
	"fmt"
	"github.com/KarolisL/lightkeeper/pkg/daemon/config"
	"github.com/KarolisL/lightkeeper/pkg/plugins/input"
	"github.com/KarolisL/lightkeeper/pkg/plugins/output"
)

type Daemon struct {
	inputs   map[string]input.Input
	outputs  map[string]output.Output
	mappings []config.Mapping
}

func (d *Daemon) Start() {
	for _, mapping := range d.mappings {
		go d.connectSync(mapping)
	}
}

func (d *Daemon) connectSync(mapping config.Mapping) {
	src := d.inputs[mapping.From].Ch()
	dest := d.outputs[mapping.To].Ch()

	for message := range src {
		dest <- message
	}
}

func NewDaemon(config *config.Config, inputMaker input.Maker, outputMaker output.OutputMaker) (*Daemon, error) {
	inputs := make(map[string]input.Input)
	for name, inputConfig := range config.Inputs {
		newInput, err := inputMaker.NewInput(inputConfig.Type, inputConfig.Params)
		if err != nil {
			return nil, fmt.Errorf("creating input %q: %w", name, err)
		}
		inputs[name] = newInput
	}

	outputs := make(map[string]output.Output)
	for name, outputConfig := range config.Outputs {
		newOutput, err := outputMaker.NewOutput(outputConfig.Type, outputConfig.Params)
		if err != nil {
			return nil, fmt.Errorf("creating output %q: %w", name, err)
		}
		outputs[name] = newOutput
	}

	mappings := config.Mappings

	d := &Daemon{inputs: inputs, outputs: outputs, mappings: mappings}

	return d, nil
}

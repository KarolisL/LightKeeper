package daemon

import (
	"fmt"
	"github.com/KarolisL/lightkeeper/pkg/common"
	"github.com/KarolisL/lightkeeper/pkg/daemon/config"
	"github.com/KarolisL/lightkeeper/pkg/plugins/input"
	"github.com/KarolisL/lightkeeper/pkg/plugins/output"
	"regexp"
)

type Daemon struct {
	inputs   map[string]input.FanoutInput
	outputs  map[string]output.Output
	mappings []config.Mapping
}

func NewDaemon(config *config.Config, inputMaker input.Maker, outputMaker output.Maker) (*Daemon, error) {
	inputs := make(map[string]input.FanoutInput)
	for name, inputConfig := range config.Inputs {
		newInput, err := inputMaker.NewFanOutInput(inputConfig.Type, inputConfig.Params)
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

func (d *Daemon) Start() {
	for i, mapping := range d.mappings {
		consume := d.makeConsumer(i, mapping)
		go consume()
	}

	for _, inp := range d.inputs {
		go func(inp input.FanoutInput) {
			err := inp.Start()
			if err != nil {
				panic(err)
			}
		}(inp)
	}
}

func (d *Daemon) makeConsumer(i int, mapping config.Mapping) func() {
	src, _ := d.inputs[mapping.From].StartListener(fmt.Sprintf("mapping#%d[%s -> %s]", i+1, mapping.From, mapping.To))
	dest := d.outputs[mapping.To].Ch()
	filters := constructFilters(mapping)

	return func() {
		for message := range src {
			if matchesAll(filters, message) {
				dest <- message
			}
		}
	}
}

func constructFilters(mapping config.Mapping) []*regexp.Regexp {
	var filters []*regexp.Regexp
	for _, filter := range mapping.Filters {
		if filter["type"] == "syslog-ng" {
			program := filter["program"]
			pattern := syslogNgProgramRegex(program)
			filters = append(filters, regexp.MustCompile(pattern))
		}
	}
	return filters
}

func matchesAll(filters []*regexp.Regexp, message common.Message) bool {
	for _, filter := range filters {
		if !filter.MatchString(message.String()) {
			return false
		}
	}

	return true
}

func syslogNgProgramRegex(program string) string {
	return fmt.Sprintf(`(?P<month>\w+)  ?(?P<day>\d+) (?P<time>\d{2}:\d{2}:\d{2}) (?P<hostname>.+) (?P<program>%s)(\[(?P<pid>\w+)\])?: (?P<msg>.*)`, program)
}

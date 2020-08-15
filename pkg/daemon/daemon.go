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
	filters := constructFilters(mapping)

	for message := range src {
		if matchesAll(filters, message) {
			dest <- message
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

func syslogNgProgramRegex(program string) string {
	return fmt.Sprintf(`(?P<month>\w+)  ?(?P<day>\d+) (?P<time>\d{2}:\d{2}:\d{2}) (?P<hostname>.+) (?P<program>%s)(\[(?P<pid>\w+)\])?: (?P<msg>.*)`, program)
}

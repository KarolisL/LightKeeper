package output

import (
	"errors"
	"fmt"
	"github.com/KarolisL/lightkeeper/pkg/daemon/config"
)

var ErrPluginTypeNotFound = errors.New("plugin not found")
var Registry = &PluginRegistry{}

type NewOutputFunc func(params config.Params) (Output, error)

type PluginRegistry struct {
	registry map[string]NewOutputFunc
}

func (r *PluginRegistry) Register(outputType string, newOutputFunc NewOutputFunc) {
	if r.registry == nil {
		r.registry = make(map[string]NewOutputFunc)
	}

	r.registry[outputType] = newOutputFunc
}

func (r *PluginRegistry) NewOutput(outputType string, params config.Params) (Output, error) {
	f, found := r.registry[outputType]
	if !found {
		return nil, fmt.Errorf("unknown output plugin type %q: %w", outputType, ErrPluginTypeNotFound)
	}

	output, err := f(params)
	if err != nil {
		return nil, fmt.Errorf("plugin %q usage error: %w", outputType, err)
	}
	return output, nil
}

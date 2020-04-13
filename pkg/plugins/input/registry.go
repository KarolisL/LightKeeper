package input

import (
	"errors"
	"fmt"
	"github.com/KarolisL/lightkeeper/pkg/daemon/config"
)

var ErrPluginTypeNotFound = errors.New("plugin not found")

var Registry = &PluginRegistry{}

type NewInputFunc func(params config.Params) (Input, error)

type PluginRegistry struct {
	registry map[string]NewInputFunc
}

func (r *PluginRegistry) Register(inputType string, newInputFunc NewInputFunc) {
	if r.registry == nil {
		r.registry = make(map[string]NewInputFunc)
	}

	r.registry[inputType] = newInputFunc
}

func (r *PluginRegistry) NewInput(inputType string, params config.Params) (Input, error) {
	f, found := r.registry[inputType]
	if !found {
		return nil, fmt.Errorf("unknown input plugin type %q: %w", inputType, ErrPluginTypeNotFound)
	}

	input, err := f(params)
	if err != nil {
		return nil, fmt.Errorf("plugin %q usage error: %w", inputType, err)
	}
	return input, nil
}

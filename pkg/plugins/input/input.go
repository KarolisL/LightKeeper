//go:generate mockgen -destination ./mock/maker.gen.go . Maker

package input

import (
	"fmt"
	"github.com/KarolisL/lightkeeper/pkg/common"
	"github.com/KarolisL/lightkeeper/pkg/daemon/config"
)

var (
	ErrAlreadyExists = fmt.Errorf("subscriber with this name already exists")
	ErrNotFound      = fmt.Errorf("subscriber with this name is not found")

	_ FanoutInput = (*fanoutInput)(nil)
)

type Input interface {
	Ch() <-chan common.Message
}

type Maker interface {
	NewInput(inputType string, params config.Params) (Input, error)
	NewFanOutInput(inputType string, params config.Params) (FanoutInput, error)
}

type FanoutInput interface {
	StartListener(name string) (<-chan common.Message, error)
	StopListener(name string) error
	Start() error
}

type fanoutInput struct {
	orig      Input
	listeners map[string]chan common.Message
}

func (e *fanoutInput) StartListener(name string) (<-chan common.Message, error) {
	if e.listeners == nil {
		e.listeners = make(map[string]chan common.Message)
	}

	if _, exists := e.listeners[name]; exists {
		return nil, ErrAlreadyExists
	}

	c := make(chan common.Message, 500)
	e.listeners[name] = c

	return c, nil
}

func (e *fanoutInput) StopListener(name string) error {
	c, found := e.listeners[name]
	if !found {
		return ErrNotFound
	}

	close(c)
	delete(e.listeners, name)

	return nil
}

func (e *fanoutInput) Start() error {
	defer e.closeAll()

	for message := range e.orig.Ch() {
		for _, c := range e.listeners {
			c <- message
		}
	}

	return nil
}

func (e *fanoutInput) closeAll() {
	for _, c := range e.listeners {
		close(c)
	}
}

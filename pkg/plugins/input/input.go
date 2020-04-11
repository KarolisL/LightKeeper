package input

import (
	"github.com/KarolisL/lightkeeper/pkg/common"
	"github.com/KarolisL/lightkeeper/pkg/daemon/config"
)

type Input interface {
	Ch() <-chan common.Message
}

type InputMaker interface {
	NewInput(inputType string, params config.Params) Input
}
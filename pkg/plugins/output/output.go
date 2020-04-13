package output

import (
	"github.com/KarolisL/lightkeeper/pkg/common"
	"github.com/KarolisL/lightkeeper/pkg/daemon/config"
)

type Output interface {
	Ch() chan<- common.Message
}

type OutputMaker interface {
	NewOutput(outputType string, params config.Params) (Output, error)
}

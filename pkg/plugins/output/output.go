//go:generate mockgen -destination ./mock/maker.gen.go . Maker
package output

import (
	"github.com/KarolisL/lightkeeper/pkg/common"
	"github.com/KarolisL/lightkeeper/pkg/daemon/config"
)

type Output interface {
	Ch() chan<- common.Message
}

type Maker interface {
	NewOutput(outputType string, params config.Params) (Output, error)
}

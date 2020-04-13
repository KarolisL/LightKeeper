package file

import (
	"github.com/KarolisL/lightkeeper/pkg/common"
	"github.com/KarolisL/lightkeeper/pkg/daemon/config"
	"github.com/KarolisL/lightkeeper/pkg/plugins/input"
	"github.com/hpcloud/tail"
)

func init() {
	input.Registry.Register("file", NewFileInput)
}

type Input struct {
	filepath string
	tail     *tail.Tail
	ch       chan common.Message
}

func (f *Input) Ch() <-chan common.Message {
	if f.ch == nil {
		f.initializeTail()
		go lineToMessageSync(f.tail.Lines, f.ch)
	}

	return f.ch
}

func (f *Input) initializeTail() {
	t, _ := tail.TailFile(f.filepath, tail.Config{ReOpen: true, Follow: true})
	f.tail = t
	f.ch = make(chan common.Message)
}

func lineToMessageSync(in chan *tail.Line, out chan common.Message) {
	for line := range in {
		out <- common.Message(line.Text)
	}
}

func NewFileInput(params config.Params) (input.Input, error) {
	inp := &Input{filepath: params["path"]}

	return inp, nil
}

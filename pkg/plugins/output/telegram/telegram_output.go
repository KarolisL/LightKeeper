package telegram

import (
	"github.com/KarolisL/lightkeeper/pkg/common"
	"github.com/KarolisL/lightkeeper/pkg/daemon/config"
	"github.com/KarolisL/lightkeeper/pkg/plugins/output"
	"github.com/yanzay/tbot/v2"
	"net/http"
)

type Output struct {
	client *tbot.Client
	chatId string

	ch chan common.Message
}

func (o *Output) Ch() chan<- common.Message {
	if o.ch == nil {
		o.ch = make(chan common.Message)
		go o.sendToTelegramSync()
	}

	return o.ch
}

func (o *Output) sendToTelegramSync() {
	for message := range o.ch {
		o.client.SendMessage(o.chatId, string(message))
	}
}

func init() {
	output.Registry.Register("telegram", NewTelegramOutput)
}

func NewTelegramOutput(params config.Params) (output.Output, error) {
	token := params["token"]
	chatId := params["chatId"]
	baseUrl := params["endpoint"]
	if baseUrl == "" {
		baseUrl = "https://api.telegram.org"
	}

	client := tbot.NewClient(token, http.DefaultClient, baseUrl)
	outp := &Output{client: client, chatId: chatId}

	return outp, nil
}

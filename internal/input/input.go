package input

import (
	"context"
	"github.com/obnahsgnaw/protocolhandler/application/register"
	"github.com/obnahsgnaw/sockethandler/service/action"
)

func init() {
	register.Register(func(p *register.Provider) {
		p.ListenRawInput(inputTransfer)
	})
}

func inputTransfer(ctx context.Context, rq *action.HandlerReq, rawInput []byte) (interAction *register.InterAction, rawOut []byte, err error) {
	rawOut = []byte("not implemented")
	return
}

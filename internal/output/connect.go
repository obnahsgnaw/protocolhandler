package output

import (
	"context"
	"errors"
	"github.com/obnahsgnaw/protocolhandler/application/register"
	"github.com/obnahsgnaw/sockethandler/service/action"
	"github.com/obnahsgnaw/socketutil/codec"
)

func init() {
	register.Register(func(p *register.Provider) {
		p.ListenOutput(connectResponse, func() codec.DataPtr {
			return nil
		}, connectResponseTransfer)
	})
}

func connectResponseTransfer(ctx context.Context, rq *action.HandlerReq, actData codec.DataPtr) (rawOut []byte, err error) {
	err = errors.New("not implemented")
	return
}

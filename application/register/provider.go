package register

import (
	"github.com/obnahsgnaw/sockethandler"
	"github.com/obnahsgnaw/sockethandler/service/action"
	"github.com/obnahsgnaw/sockethandler/service/proto/impl"
	"github.com/obnahsgnaw/socketutil/codec"
	"go.uber.org/zap"
)

type Provider struct {
	s *sockethandler.Handler
}

var _provider *Provider

func init() {
	_provider = &Provider{}
}

func (d *Provider) ListenRawInput(transfer InputTransfer) {
	_dispatcher.listenRawInput(transfer)
}

func (d *Provider) ListenOutput(actId codec.ActionId, structure action.DataStructure, handler OutputHandler) {
	_dispatcher.listenOutput(actId, structure, handler)
}

func (d *Provider) Logger() *zap.Logger {
	return d.s.Logger()
}

func (d *Provider) TcpGateway() *impl.Gateway {
	return d.s.TcpGateway()
}

func (d *Provider) WssGateway() *impl.Gateway {
	return d.s.WssGateway()
}

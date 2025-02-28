package application

import (
	"github.com/obnahsgnaw/application"
	"github.com/obnahsgnaw/application/endtype"
	"github.com/obnahsgnaw/application/pkg/url"
	"github.com/obnahsgnaw/protocolhandler/application/register"
	"github.com/obnahsgnaw/protocolhandler/config"
	"github.com/obnahsgnaw/rpc"
	"github.com/obnahsgnaw/sockethandler"
	"github.com/obnahsgnaw/sockethandler/sockettype"
)

var (
	rpsIns  *sockethandler.ManagedRpc
	rpsIp   string
	rpsPort int
)

func New(app *application.Application, o ...Option) *sockethandler.Handler {
	validate()
	for _, opt := range o {
		opt()
	}
	if rpsIns == nil {
		rpsIns = newRps(app, rpsIp, rpsPort)
	}
	h := sockethandler.New(app, rpsIns, "device", "protocol:"+config.Name, config.Name, endtype.Frontend, sockettype.TCP)
	register.Load(h, config.ModelNames, config.ActionId)
	app.AddServer(h)
	return h
}

func newRps(app *application.Application, ip string, port int) *sockethandler.ManagedRpc {
	l, _ := rpc.NewListener(url.Host{Ip: ip, Port: port})
	rps := sockethandler.NewRpc(app, "device", "protocol:"+config.Name, endtype.Frontend, sockettype.TCP, l, nil, rpc.RegEnable())
	return rps
}

func validate() {
	if config.Name == "" {
		panic("application config name is required")
	}
	if config.ModelNames == "" {
		panic("application config models is required")
	}
}

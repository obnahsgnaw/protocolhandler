package application

import (
	"context"
	"github.com/obnahsgnaw/application"
	"github.com/obnahsgnaw/application/endtype"
	"github.com/obnahsgnaw/application/pkg/url"
	"github.com/obnahsgnaw/protocolhandler/application/register"
	"github.com/obnahsgnaw/protocolhandler/config"
	_ "github.com/obnahsgnaw/protocolhandler/internal"
	"github.com/obnahsgnaw/rpc"
	handlerv1 "github.com/obnahsgnaw/socketapi/gen/handler/v1"
	"github.com/obnahsgnaw/socketgateway"
	"github.com/obnahsgnaw/socketgateway/pkg/socket"
	"github.com/obnahsgnaw/sockethandler"
	"github.com/obnahsgnaw/sockethandler/service/action"
	"github.com/obnahsgnaw/socketutil/codec"
	"strconv"
	"time"
)

var (
	rpsIns  *sockethandler.ManagedRpc
	rpsIp   string
	rpsPort int
)

func BindToGateway(s *socketgateway.Server) {
	validate()
	act := codec.Action{Id: config.ActionId, Name: "raw:" + config.ModelNames}
	s.ActionManager().RegisterHandlerAction(act, func() codec.DataPtr {
		return &handlerv1.RawRequest{}
	}, func(c socket.Conn, data codec.DataPtr) (respAction codec.Action, respData codec.DataPtr) {
		var err error
		rq := data.(*handlerv1.RawRequest)
		resp := &handlerv1.RawResponse{}
		data = resp

		ctx, cl := context.WithTimeout(context.Background(), time.Second)
		defer cl()

		var u *action.User
		if c.Context().Authed() {
			u = &action.User{
				Id:   uint32(int(c.Context().User().Id)),
				Name: c.Context().User().Name,
				Attr: c.Context().User().Attr,
			}
			if u.Attr == nil {
				u.Attr = make(map[string]string)
			}
		}
		var target *action.Target
		if c.Context().Authentication() != nil {
			target = &action.Target{
				Type:     c.Context().Authentication().Type,
				Id:       c.Context().Authentication().Id,
				Cid:      c.Context().Authentication().Cid,
				Uid:      c.Context().Authentication().Uid,
				Protocol: c.Context().Authentication().Protocol,
			}
		}
		req := action.NewHandlerReq(s.Rpc().Host().String(), act, int64(c.Fd()), u, data, c.Context().IdMap(), target)

		if rq.ActionId == 0 {
			respAction, resp.Data, resp.SubActions, err = register.Dispatcher().DispatchInput(ctx, req, rq.Data)
			if err != nil {
				s.Logger().Error("transfer input failed, err=" + err.Error())
			} else {
				s.Logger().Debug("transfer input:" + string(rq.Data) + ",out:action=" + respAction.String())
			}
		} else {
			resp.Data, err = register.Dispatcher().DispatchOutput(ctx, req, codec.ActionId(rq.ActionId), rq.Data)
			if err != nil {
				s.Logger().Error("transfer output failed, err=" + err.Error())
			} else {
				s.Logger().Debug("transfer output:" + string(rq.Data) + ",out:action=" + strconv.Itoa(int(rq.ActionId)))
			}
		}
		return
	})
}

func BindToHandler(s *sockethandler.Handler) {
	validate()
	register.Load(s, config.ModelNames, config.ActionId)
}

func New(app *application.Application, o ...Option) *sockethandler.Handler {
	validate()
	for _, opt := range o {
		opt()
	}
	if rpsIns == nil {
		rpsIns = newRps(app, rpsIp, rpsPort)
	}
	h := sockethandler.New(app, rpsIns, "device", "protocol:"+config.Name, config.Name, endtype.Frontend, "outer")
	BindToHandler(h)
	app.AddServer(h)
	return h
}

func newRps(app *application.Application, ip string, port int) *sockethandler.ManagedRpc {
	l, _ := rpc.NewListener(url.Host{Ip: ip, Port: port})
	rps := sockethandler.NewRpc(app, "device", "protocol:"+config.Name, endtype.Frontend, "outer", l, nil, rpc.RegEnable())
	return rps
}

func validate() {
	if config.Name == "" {
		panic("protocol handler: application config name is required")
	}
	if config.ModelNames == "" {
		panic("protocol handler: application config models is required")
	}
	if config.ActionId == 0 {
		panic("protocol handler: application config action id is required")
	}
}

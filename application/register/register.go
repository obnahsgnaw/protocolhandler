package register

import (
	"context"
	handlerv1 "github.com/obnahsgnaw/socketapi/gen/handler/v1"
	"github.com/obnahsgnaw/sockethandler"
	"github.com/obnahsgnaw/sockethandler/service/action"
	"github.com/obnahsgnaw/socketutil/codec"
	"strconv"
	"strings"
)

func Load(s *sockethandler.Handler, modelNames string, actionId int) {
	_provider.s = s
	for _, modelName := range strings.Split(modelNames, ",") {
		if modelName != "" {
			s.Listen(codec.Action{
				Id:   codec.ActionId(actionId),
				Name: "raw:" + modelName,
			}, func() codec.DataPtr {
				return &handlerv1.RawRequest{}
			}, func(ctx context.Context, req *action.HandlerReq) (respAct codec.Action, data codec.DataPtr, err error) {
				rq := req.Data.(*handlerv1.RawRequest)
				resp := &handlerv1.RawResponse{}
				data = resp
				if rq.ActionId == 0 {
					respAct, resp.Data, err = _dispatcher.DispatchInput(ctx, req, rq.Data)
					if err != nil {
						s.Logger().Error("transfer input failed, err=" + err.Error())
					} else {
						s.Logger().Debug("transfer input:" + string(rq.Data) + ",out:action=" + respAct.String())
					}
				} else {
					resp.Data, err = _dispatcher.DispatchOutput(ctx, req, codec.ActionId(rq.ActionId), rq.Data)
					if err != nil {
						s.Logger().Error("transfer output failed, err=" + err.Error())
					} else {
						s.Logger().Debug("transfer output:" + string(rq.Data) + ",out:action=" + strconv.Itoa(int(rq.ActionId)))
					}
				}
				return
			})
		}
	}
}

func Register(cb func(*Provider)) {
	cb(_provider)
}

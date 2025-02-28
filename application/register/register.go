package register

import (
	"context"
	handlerv1 "github.com/obnahsgnaw/socketapi/gen/handler/v1"
	"github.com/obnahsgnaw/sockethandler"
	"github.com/obnahsgnaw/sockethandler/service/action"
	"github.com/obnahsgnaw/socketutil/codec"
	"strings"
)

func Load(s *sockethandler.Handler, modelNames string, actionId int) {
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
					respAct, resp.Data, err = Dispatcher().dispatchInput(rq.Data)
				} else {
					resp.Data, err = Dispatcher().dispatchOutput(codec.ActionId(rq.ActionId), rq.Data)
				}
				return
			})
		}
	}
}

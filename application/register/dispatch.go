package register

import (
	"context"
	"errors"
	handlerv1 "github.com/obnahsgnaw/socketapi/gen/handler/v1"
	"github.com/obnahsgnaw/sockethandler/service/action"
	"github.com/obnahsgnaw/socketutil/codec"
	"google.golang.org/protobuf/proto"
	"strconv"
)

type Dispatch struct {
	dataBuilder    codec.DataBuilder
	inputTransfer  InputTransfer
	outputTransfer map[codec.ActionId]outHandler
}

type InterAction struct {
	Action     codec.Action
	Data       proto.Message
	SubActions []*SubAction
}
type SubAction struct {
	target   string
	ActionId codec.ActionId
	Data     proto.Message
}
type outHandler struct {
	structure action.DataStructure
	handler   OutputHandler
}
type InputTransfer func(ctx context.Context, rq *action.HandlerReq, rawInput []byte) (interAction *InterAction, rawOut []byte, err error)
type OutputHandler func(ctx context.Context, rq *action.HandlerReq, actData codec.DataPtr) (rawOut []byte, err error)

var _dispatcher *Dispatch

func init() {
	_dispatcher = &Dispatch{
		dataBuilder: codec.NewProtobufDataBuilder(),
		inputTransfer: func(ctx context.Context, rq *action.HandlerReq, rawInput []byte) (interAction *InterAction, rawOut []byte, err error) {
			err = errors.New("raw inout package transfer not implemented")
			return
		},
		outputTransfer: make(map[codec.ActionId]outHandler),
	}
}

func Dispatcher() *Dispatch {
	return _dispatcher
}

// DispatchInput act=0 则直接响应， 否则会转发给实际的action handler
func (d *Dispatch) DispatchInput(ctx context.Context, rq *action.HandlerReq, rawIn []byte) (act codec.Action, actData []byte, subActions []*handlerv1.SubAction, err error) {
	interAction, rawOut, transErr := d.inputTransfer(ctx, rq, rawIn)
	if transErr != nil {
		err = transErr
		return
	}
	if interAction == nil {
		act = codec.NewAction(0, "")
		actData = rawOut
		return
	}
	act = interAction.Action
	actData, err = d.dataBuilder.Pack(interAction.Data)
	if err != nil {
		err = errors.New("pack data failed: " + err.Error())
		return
	}
	for _, subAction := range interAction.SubActions {
		subData, subErr := d.dataBuilder.Pack(subAction.Data)
		if subErr != nil {
			err = errors.New("sub action[" + strconv.Itoa(int(subAction.ActionId)) + "] pack data failed: " + subErr.Error())
			return
		}
		subActions = append(subActions, &handlerv1.SubAction{
			ActionId: uint32(subAction.ActionId),
			Data:     subData,
			Target:   subAction.target,
		})
	}
	return
}

// DispatchOutput 转化实际action handler的数据成原始数据
func (d *Dispatch) DispatchOutput(ctx context.Context, rq *action.HandlerReq, act codec.ActionId, actData []byte) (rawOut []byte, err error) {
	if handler, ok := d.outputTransfer[act]; ok {
		ptr := handler.structure()
		if err = d.dataBuilder.Unpack(actData, ptr); err != nil {
			return
		}
		rawOut, err = handler.handler(ctx, rq, ptr)
		return
	}
	err = errors.New("action transfer not supported")
	return
}

func (d *Dispatch) listenRawInput(transfer InputTransfer) {
	if transfer != nil {
		d.inputTransfer = transfer
	}
}

func (d *Dispatch) listenOutput(actId codec.ActionId, structure action.DataStructure, handler OutputHandler) {
	d.outputTransfer[actId] = outHandler{
		structure: structure,
		handler:   handler,
	}
}

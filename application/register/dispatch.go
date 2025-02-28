package register

import "github.com/obnahsgnaw/socketutil/codec"

type Dispatch struct{}

var _dispatcher *Dispatch

func init() {
	_dispatcher = &Dispatch{}
}

func Dispatcher() *Dispatch {
	return _dispatcher
}

// act=0 则直接响应， 否则会转发给实际的action handler
func (d *Dispatch) dispatchInput(rawIn []byte) (act codec.Action, actData []byte, err error) {
	actData = []byte(`ok`)
	return
}

// 转化实际action handler的数据成原始数据
func (d *Dispatch) dispatchOutput(act codec.ActionId, actData []byte) (rawOut []byte, err error) {
	return
}

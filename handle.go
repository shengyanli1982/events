package events

// handleFuncs 是一个结构体，它包含两个字段：origFunc 和 wrapFunc，都是 MessageHandleFunc 类型的函数。
// handleFuncs is a structure that contains two fields: origFunc and wrapFunc, both of which are functions of type MessageHandleFunc.
type handleFuncs struct {
	// origFunc 是原始的消息处理函数。
	// origFunc is the original message handling function.
	origFunc MessageHandleFunc

	// wrapFunc 是包装后的消息处理函数。
	// wrapFunc is the wrapped message handling function.
	wrapFunc MessageHandleFunc
}

// newHandleFuncs 是一个函数，它返回一个新的 handleFuncs 实例。
// newHandleFuncs is a function that returns a new instance of handleFuncs.
func newHandleFuncs() *handleFuncs {
	return &handleFuncs{}
}

// SetOrigMsgHandleFunc 是 handleFuncs 的一个方法，它设置 origFunc 字段的值。
// SetOrigMsgHandleFunc is a method of handleFuncs that sets the value of the origFunc field.
func (h *handleFuncs) SetOrigMsgHandleFunc(fn MessageHandleFunc) {
	h.origFunc = fn
}

// GetOrigMsgHandleFunc 是 handleFuncs 的一个方法，它返回 origFunc 字段的值。
// GetOrigMsgHandleFunc is a method of handleFuncs that returns the value of the origFunc field.
func (h *handleFuncs) GetOrigMsgHandleFunc() MessageHandleFunc {
	return h.origFunc
}

// SetWrapMsgHandleFunc 是 handleFuncs 的一个方法，它设置 wrapFunc 字段的值。
// SetWrapMsgHandleFunc is a method of handleFuncs that sets the value of the wrapFunc field.
func (h *handleFuncs) SetWrapMsgHandleFunc(fn MessageHandleFunc) {
	h.wrapFunc = fn
}

// GetWrapMsgHandleFunc 是 handleFuncs 的一个方法，它返回 wrapFunc 字段的值。
// GetWrapMsgHandleFunc is a method of handleFuncs that returns the value of the wrapFunc field.
func (h *handleFuncs) GetWrapMsgHandleFunc() MessageHandleFunc {
	return h.wrapFunc
}

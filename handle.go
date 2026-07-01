package events

// handleFuncs 保存每个 topic 的处理函数元信息：原始函数、包装后函数及是否 once 语义。
type handleFuncs struct {
	origFunc MessageHandleFunc
	wrapFunc MessageHandleFunc
	once     bool
}

func newHandleFuncs() *handleFuncs {
	return &handleFuncs{}
}

func (h *handleFuncs) SetOrigMsgHandleFunc(fn MessageHandleFunc) {
	h.origFunc = fn
}

func (h *handleFuncs) GetOrigMsgHandleFunc() MessageHandleFunc {
	return h.origFunc
}

func (h *handleFuncs) SetWrapMsgHandleFunc(fn MessageHandleFunc) {
	h.wrapFunc = fn
}

func (h *handleFuncs) GetWrapMsgHandleFunc() MessageHandleFunc {
	return h.wrapFunc
}

func (h *handleFuncs) SetOnce(once bool) {
	h.once = once
}

func (h *handleFuncs) IsOnce() bool {
	return h.once
}

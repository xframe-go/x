package xevent

type FuncListener func(payload any) error

func (fn FuncListener) Handle(payload any) error {
	return fn(payload)
}

const defaultBufferSize int = 32

type AsyncFuncListener struct {
	fn     FuncListener
	buffer int
}

func NewAsyncFuncListener(fn FuncListener, buffer ...int) *AsyncFuncListener {
	size := defaultBufferSize
	if len(buffer) > 0 {
		size = buffer[0]
	}

	return &AsyncFuncListener{
		fn:     fn,
		buffer: size,
	}
}

func (fn AsyncFuncListener) Handle(payload any) error {
	return fn.fn(payload)
}

func (fn AsyncFuncListener) AsyncBuffer() int {
	return fn.buffer
}

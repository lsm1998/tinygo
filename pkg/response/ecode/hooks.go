package ecode

type ErrHandle func(err error)

var (
	errHook ErrHandle
)

// InitErrHook 这个是全局的，请勿在程序中调用多次
func InitErrHook(hook ErrHandle) {
	errHook = hook
}

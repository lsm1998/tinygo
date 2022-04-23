package tinygo

import "os"

const (
	defaultAppMod   = "local"
	defaultAppTrace = "false"
)

var (
	AppMod   string
	AppTrace string
)

func init() {
	AppMod = os.Getenv("APP_MOD")
	if AppMod == "" {
		AppMod = defaultAppMod
	}
	AppTrace = os.Getenv("APP_TRACE")
	if AppTrace == "" {
		AppTrace = defaultAppTrace
	}
}

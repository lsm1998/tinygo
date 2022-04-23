package response

import (
	"fmt"
	"github.com/lsm1998/tinygo/pkg/response/ecode"
	"github.com/pkg/errors"
)

func RpcError(err error) error {
	if e, ok := err.(*ecode.Errno); ok {
		return errors.New(fmt.Sprintf("%s<sp>%d", e.Msg, e.Code))
	} else {
		return err
	}
}

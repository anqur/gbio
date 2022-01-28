package codec

import (
	"fmt"

	"github.com/anqur/gbio/base"
)

var (
	ErrBadMsgTag  = fmt.Errorf("%w: unknown message tag", base.Err)
	ErrBadMsgType = fmt.Errorf("%w: unknown message type", base.Err)
)

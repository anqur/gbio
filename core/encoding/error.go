package encoding

import (
	"fmt"

	"github.com/anqur/gbio/core/gbioerr"
)

var (
	ErrBadMsgTag  = fmt.Errorf("%w: unknown message tag", gbioerr.Err)
	ErrBadMsgType = fmt.Errorf("%w: unknown message type", gbioerr.Err)
)

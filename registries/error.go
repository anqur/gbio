package registries

import (
	"fmt"

	"github.com/anqur/gbio/base"
)

var (
	ErrEndpointNotFound = fmt.Errorf("%w: endpoint not found", base.Err)
	ErrEmptyEndpoints   = fmt.Errorf("%w: unexpected empty endpoints", base.Err)
)

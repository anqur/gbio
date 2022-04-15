package registries

import (
	"fmt"
	"github.com/anqur/gbio/pkg/gbioerr"
)

var (
	ErrEndpointNotFound = fmt.Errorf("%w: endpoint not found", gbioerr.Err)
	ErrEmptyEndpoints   = fmt.Errorf("%w: unexpected empty endpoints", gbioerr.Err)
)

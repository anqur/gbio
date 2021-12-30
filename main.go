package gbio

import (
	"log"

	"github.com/anqur/gbio/internal/errors"
	"github.com/anqur/gbio/internal/loggers"
)

var Err = errors.Err

func UseErrorLogger(l *log.Logger) { loggers.Error = l }

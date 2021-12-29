package gbio

import (
	"log"

	"github.com/anqur/gbio/internal/errors"
	"github.com/anqur/gbio/internal/loggers"
)

var Err = errors.Err

func UseInfoLogger(l *log.Logger)  { loggers.Info = l }
func UseErrorLogger(l *log.Logger) { loggers.Error = l }

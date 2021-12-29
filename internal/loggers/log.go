package loggers

import (
	"log"
	"os"
)

var Info = log.New(
	os.Stderr,
	"INFO ",
	log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix,
)

var Error = log.New(
	os.Stderr,
	"ERROR ",
	log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix,
)

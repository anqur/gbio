package loggers

import (
	"log"
	"os"
)

const DefaultFlag = log.LstdFlags | log.Lmicroseconds | log.Lmsgprefix

var Error = log.New(os.Stderr, "ERROR ", DefaultFlag)

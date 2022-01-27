package logging

import (
	"log"
	"os"
)

const DefaultFlag = log.LstdFlags | log.Lmicroseconds | log.Lmsgprefix

var (
	Info  = log.New(os.Stderr, "INFO ", DefaultFlag)
	Error = log.New(os.Stderr, "ERROR ", DefaultFlag)
)

package logging

import (
	"log"
	"os"
)

type LogLiner interface{ Println(a ...any) }

const DefaultFlag = log.LstdFlags | log.Lmicroseconds | log.Lmsgprefix

var (
	Info  LogLiner = log.New(os.Stderr, "INFO ", DefaultFlag)
	Error LogLiner = log.New(os.Stderr, "ERROR ", DefaultFlag)
)

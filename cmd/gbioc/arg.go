package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/anqur/gbio/internal/utils"
)

var (
	isVersion     bool
	file          string
	outDir        string
	target        string
	marshaller    string
	discriminator string
	ctxKey        string
)

const (
	TargetServerGo = "server-go"
	TargetClientGo = "client-go"
)

var targets = []string{
	TargetServerGo,
	TargetClientGo,
}

const (
	MarshallerJSON = "json"
)

var marshallers = []string{
	MarshallerJSON,
}

func oneOfHelp(enums []string, help string) string {
	return fmt.Sprintf("%s, e.g. %s", help, strings.Join(enums, ", "))
}

func parseArgs() {
	flag.BoolVar(
		&isVersion,
		"v",
		false,
		"show version",
	)
	flag.StringVar(
		&file,
		"i",
		"",
		"input file",
	)
	flag.StringVar(
		&outDir,
		"o",
		".",
		"output directory",
	)
	flag.StringVar(
		&target,
		"t",
		TargetServerGo,
		oneOfHelp(targets, "code generation target"),
	)
	flag.StringVar(
		&marshaller,
		"m",
		MarshallerJSON,
		oneOfHelp(marshallers, "marshalling format"),
	)
	flag.StringVar(
		&discriminator,
		"enumtag",
		"_t",
		"discriminator key for enum types",
	)
	flag.StringVar(
		&ctxKey,
		"ctxkey",
		"ctx",
		"context key for (un)marshalling context",
	)

	flag.Usage = showUsage
	flag.Parse()
}

func invalidArg(format string, v ...interface{}) {
	showUsage()
	log.Fatalf(format, v...)
}

func validateArgs() {
	if !utils.OneOf(target, targets) {
		invalidArg("Invalid code generation target: %s\n", target)
	}
	if !utils.OneOf(marshaller, marshallers) {
		invalidArg("Invalid marshalling format: %s\n", marshaller)
	}
}

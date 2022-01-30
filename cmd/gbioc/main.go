package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/anqur/gbio/internal/codegens"
	"github.com/anqur/gbio/internal/parsers"
	"github.com/anqur/gbio/internal/utils"
)

var (
	isVersion bool

	file          string
	outDir        string
	target        string
	marshaller    string
	discriminator string
)

const (
	Version = "0.1.0"

	TargetServerGo = "server-go"
	TargetClientGo = "client-go"

	MarshallerJSON = "json"

	DefaultEnumTag = "_t"
)

var (
	targets = []string{
		TargetServerGo,
		TargetClientGo,
	}

	marshallers = []string{
		MarshallerJSON,
	}
)

func oneOfHelp(enums []string, help string) string {
	return fmt.Sprintf("%s, e.g. %s", help, strings.Join(enums, ", "))
}

func invalidArgs(format string, args ...interface{}) {
	showUsage()
	log.Fatalln(fmt.Sprintf(format, args...))
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
		DefaultEnumTag,
		"discriminator key for enum types",
	)
	flag.Usage = showUsage
	flag.Parse()
}

func validateArgs() {
	if !utils.OneOf(target, targets) {
		invalidArgs(
			"Invalid code generation target %q, expected %v",
			target,
			targets,
		)
	}
	if !utils.OneOf(marshaller, marshallers) {
		invalidArgs(
			"Invalid marshalling format: %q, expected %v",
			marshaller,
			marshallers,
		)
	}
	if file == "" {
		invalidArgs("Unexpected empty input file")
	}
}

func showUsage() {
	showVersion()
	fmt.Println("Usage:")
	flag.PrintDefaults()
}

func showVersion() {
	fmt.Println("gbioc v" + Version)
}

func main() {
	parseArgs()
	validateArgs()

	if isVersion {
		showVersion()
		return
	}

	(&parsers.Parser{
		Codegen: codegens.Codegen{
			OutDir:        outDir,
			Target:        target,
			Marshaller:    marshaller,
			Discriminator: discriminator,
		},
		File: file,
	}).Parse().Generate()
}

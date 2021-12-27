package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

const (
	version = "0.1.0"
)

var (
	isVersion bool
	file      string
	outDir    string
	target    string
)

const (
	TargetServerGo = "server-go"
	TargetClientGo = "client-go"
)

var targets = []string{
	TargetServerGo,
	TargetClientGo,
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
		fmt.Sprintf(
			"code generation target (%s)",
			strings.Join(targets, ", "),
		),
	)

	flag.Usage = showUsage
	flag.Parse()
}

func validateArgs() {
	if !stringIn(target, targets) {
		showUsage()
		log.Fatalf("Invalid code generation target: %s\n", target)
	}
}

func showUsage() {
	showVersion()
	fmt.Println("Usage:")
	flag.PrintDefaults()
}

func showVersion() {
	fmt.Println("gbioc v" + version)
}

func main() {
	parseArgs()
	validateArgs()

	if isVersion {
		showVersion()
		return
	}
}

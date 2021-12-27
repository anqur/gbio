package main

import (
	"flag"
	"fmt"
)

const version = "0.1.0"

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

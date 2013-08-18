package main

import (
	"flag"
	"github.com/fudanchii/sifr/irc"
	"os"
)

// Constants
var (
	VERSION = "v0.0.0"
)

// Flags
var (
	flVersion = flag.Bool("version", false, "Show current version then exit.")
)

func showVersion() {
	os.Stderr.WriteString(VERSION + "\n")
}

func main() {
	flag.Parse()
	if *flVersion {
		showVersion()
		os.Exit(0)
	}
}

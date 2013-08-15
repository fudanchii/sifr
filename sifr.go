package main

import (
    "flag"
    "github.com/fudanchii/sifr/irc"
	"os"
)

var (
	version = flag.Bool("version", false, "Show current version then exit.")
)

func showVersion() {
	os.Stderr.WriteString("v0.0.0\n")
}

func main() {
	flag.Parse()
	if *version {
		showVersion()
		os.Exit(0)
	}
}

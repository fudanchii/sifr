package main

import (
	"flag"
	"github.com/fudanchii/sifr/irc"
	"log"
	"os"
)

// Constants
var (
	APPNAME = "Sifr"
	VERSION = "v0.0.0"
)

// Flags
var (
	flVersion  = flag.Bool("version", false, "Show current version then exit.")
	flNick     = flag.String("nick", "shifuru", "Nickname to use.")
	flUsername = flag.String("username", "shifuru", "Username to use.")
	flRealname = flag.String("realname", "shifuru", "Realname to use.")
	flPassword = flag.String("password", "", "User's password")
	flServer   = flag.String("server", "", "Server to connect to.")
	flHelp     = flag.Bool("help", false, "Display usage, then exit")
	flDebug    = flag.Bool("debug", false, "Display debug messages.")
)

func showVersion() {
	os.Stderr.WriteString(APPNAME + "-" + VERSION + "\n")
}

func intro() {
	if *flVersion {
		showVersion()
		os.Exit(0)
	}
	if *flHelp {
		showVersion()
		flag.Usage()
		os.Exit(0)
	}
}

func main() {
	flag.Parse()
	intro()
	user := irc.NewUser(*flNick, *flUsername, *flRealname, *flPassword)
	client, err := irc.Connect(*flServer, *user)
	if err != nil {
		log.Fatal(err.Error() + "\n")
	}
	<-client.Errorchan
}

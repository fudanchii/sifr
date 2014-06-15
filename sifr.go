package main

import (
	"flag"
	"github.com/fudanchii/sifr/irc"
	"log"
	"os"
)

var (
	APPNAME = "Sifr"
	VERSION = "v0.0.0"
)

var (
	flNick     = flag.String("nick", "shifuru", "Nickname to use.")
	flUsername = flag.String("username", "shifuru", "Username to use.")
	flRealname = flag.String("realname", "shifuru", "Realname to use.")
	flPassword = flag.String("password", "", "User's password")
	flServer   = flag.String("server", "", "Server to connect to.")
	flDebug    = flag.Bool("debug", false, "Display debug messages.")
	flVersion  = flag.Bool("version", false, "Show version.")
)

func showVersion() {
	os.Stderr.WriteString(APPNAME + "-" + VERSION + "\n")
}

func main() {
	flag.Parse()
	if *flVersion {
		showVersion()
		os.Exit(0)
	}
	user := irc.NewUser(*flNick, *flUsername, *flRealname, *flPassword)
	client, err := irc.Connect(*flServer, *user)
	if err != nil {
		log.Fatalf("%s", err)
	}
	<-client.Errorchan
}

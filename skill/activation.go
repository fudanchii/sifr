package skill

import (
	"github.com/fudanchii/sifr/irc"
	"strings"
)

func nocmd(txt, cmd string, hasArg bool) bool {
	offset := len(cmd)
	if hasArg {
		offset += 2
	}
	if len(txt) < offset {
		return true
	}
	return !strings.HasPrefix(txt, cmd)
}

func ActivateFor(c *irc.Client) {
	c.AddHandler("PRIVMSG", func(msg *irc.Message) {
		flipCoin(c, msg)
	})

	c.AddHandler("PRIVMSG", func(msg *irc.Message) {
		joinChannel(c, msg)
	})
}

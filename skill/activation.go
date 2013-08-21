package skill

import (
	"github.com/fudanchii/sifr/irc"
	"strings"
)

type archetype struct {
	cmd        string
	hasArg     bool
	authorized bool
	function   func(c *irc.Client, m *irc.Message)
}

func forgeSkill(arch archetype, c *irc.Client) irc.MessageHandler {
	return func(msg *irc.Message) {
		if nocmd(msg.Body, arch.cmd, arch.hasArg) {
			return
		}
		switch {
		case arch.authorized && canPermit(msg):
			fallthrough
		case !arch.authorized:
			arch.function(c, msg)
		}
	}
}

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
	c.AddHandler("PRIVMSG", forgeSkill(archetype{
		cmd:        ".authorize",
		hasArg:     true,
		authorized: true,
		function:   addAuthorizedIdent,
	}, c))
	c.AddHandler("PRIVMSG", forgeSkill(archetype{
		cmd:        ".c",
		hasArg:     true,
		authorized: false,
		function:   flipCoin,
	}, c))
	c.AddHandler("PRIVMSG", forgeSkill(archetype{
		cmd:        ".j",
		hasArg:     true,
		authorized: true,
		function:   joinChannel,
	}, c))
	c.AddHandler("PRIVMSG", forgeSkill(archetype{
		cmd:        ".leave",
		hasArg:     true,
		authorized: true,
		function:   leaveChannel,
	}, c))
}

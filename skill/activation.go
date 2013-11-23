package skill

import (
	"github.com/fudanchii/sifr/irc"
	"strings"
)

type archetype struct {
	cmd        string
	argc       uint
	authorized bool
	function   func(c *irc.Client, m *irc.Message)
}

func forgeSkill(arch *archetype, c *irc.Client) irc.MessageHandler {
	return func(msg *irc.Message) {
		if nocmd(msg.Body, arch.cmd, arch.argc) {
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

func nocmd(txt, cmd string, argc uint) bool {
	offset := len(cmd)
	if argc > 0 {
		offset += 2
	}
	if len(txt) < offset {
		return true
	}
	return !strings.HasPrefix(txt, cmd)
}

func ActivateFor(c *irc.Client) {
	c.AddHandler("PRIVMSG", forgeSkill(&archetype{
		cmd:        ".authorize",
		argc:       1,
		authorized: true,
		function:   addAuthorizedIdent,
	}, c))
	c.AddHandler("PRIVMSG", forgeSkill(&archetype{
		cmd:        ".c",
		argc:       1,
		authorized: false,
		function:   flipCoin,
	}, c))
	c.AddHandler("PRIVMSG", forgeSkill(&archetype{
		cmd:        ".j",
		argc:       1,
		authorized: true,
		function:   joinChannel,
	}, c))
	c.AddHandler("PRIVMSG", forgeSkill(&archetype{
		cmd:        ".leave",
		argc:       1,
		authorized: true,
		function:   leaveChannel,
	}, c))
}

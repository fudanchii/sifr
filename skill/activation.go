package skill

import (
	"github.com/fudanchii/sifr/irc"
	"strings"
)

type archetype func(c *irc.Client, m *irc.Message)

type archMeta struct {
    cmd string
    hasArg bool
    authorized bool
    archFn archetype
}

func forgeSkill(archmeta archMeta, c *irc.Client) irc.MessageHandler {
	return func(msg *irc.Message) {
		if nocmd(msg.Body, archmeta.cmd, archmeta.hasArg) {
			return
		}
		archmeta.archFn(c, msg)
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
	c.AddHandler("PRIVMSG", forgeSkill(".c", true, flipCoin, c))
	c.AddHandler("PRIVMSG", forgeSkill(".j", true, joinChannel, c))
	c.AddHandler("PRIVMSG", forgeSkill(".leave", true, leaveChannel, c))
}

package skill

import (
	"github.com/fudanchii/sifr/irc"
	"strings"
)

func joinChannel(c *irc.Client, msg *irc.Message) {
	c.Join(strings.TrimSpace(msg.Body[3:]), "")
}

func leaveChannel(c *irc.Client, msg *irc.Message) {
	c.Part(strings.TrimSpace(msg.Body[7:]))
}

package irc

import (
	msg "github.com/fudanchii/sifr/irc/message"

	"strings"
)

var ctcpHandlers = map[string]MessageHandler{}

func AddCTCPHandler(cmd string, fn MessageHandler) {
	cmd = strings.ToUpper(cmd)
	ctcpHandlers[cmd] = fn
}

type defaultCTCPHandler struct{}

func (p *defaultCTCPHandler) Handle(c *Client, message *msg.Message) {
	if message.IsCTCP() && c.User.OwnThis(message) {
		p.handleCTCP(c, message)
	}
}

func (p *defaultCTCPHandler) handleCTCP(c *Client, m *msg.Message) {
	rawcmd := msg.UntagCTCP(m.Body)
	cmd := strings.Split(rawcmd, " ")[0]
	sender := m.FromNick()

	if h, exists := ctcpHandlers[cmd]; exists {
		h.Handle(c, m)
		return
	}

	switch {
	case cmd == "VERSION":
		c.ResponseCTCP(sender, "VERSION Sifr:0.0.0")
	case cmd == "SOURCE":
		c.ResponseCTCP(sender, "SOURCE https://github.com/fudanchii/sifr")
	case cmd == "PING":
		c.ResponseCTCP(sender, rawcmd)
	}
}

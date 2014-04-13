package irc

import (
	"strings"

	msg "github.com/fudanchii/sifr/irc/message"
)

// Add default handler for incoming messages.
func (c *Client) setupHandlers() {
	if c.msgHandlers == nil {
		c.msgHandlers = MessageHandlers{}
		c.AddHandler("PING", &pingHandler{})
		c.AddHandler("PRIVMSG", &privMsgDefaultHandler{})
	}
}

// Assign MessageHandler to Client's msgHandlers list.
// MessageHandler then will executed upon incoming Message,
// like a transparent filter chain.
func (c *Client) AddHandler(cmd string, fn MessageHandler) {
	cmd = strings.ToUpper(cmd)
	if _, exists := c.msgHandlers[cmd]; !exists {
		c.msgHandlers[cmd] = []MessageHandler{}
	}
	c.msgHandlers[cmd] = append(c.msgHandlers[cmd], fn)
}

// -------------------------
// Default handler for privmsg

type privMsgDefaultHandler struct{}

func (p *privMsgDefaultHandler) Handle(c *Client, m *msg.Message) {
	if m.IsCTCP() && c.User.IsMsgForMe(m) {
		p.handleCTCP(c, m)
	}
}

func (p *privMsgDefaultHandler) handleCTCP(c *Client, m *msg.Message) {
	cmd := msg.UntagCTCP(m.Body)
	sender := m.FromNick()
	switch {
	case cmd == "VERSION":
		c.ResponseCTCP(sender, "VERSION Sifr:0.0.0")
	case cmd == "SOURCE":
		c.ResponseCTCP(sender, "SOURCE https://github.com/fudanchii/sifr")
	case strings.HasPrefix(cmd, "PING "):
		c.ResponseCTCP(sender, cmd)
	}
}

// -------------------------
// ping handler

type pingHandler struct{}

func (p *pingHandler) Handle(c *Client, m *msg.Message) {
	c.Pong(m.Body)
}

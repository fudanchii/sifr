package irc

import (
	"strings"
)

// Add default handler for incoming messages.
func (c *Client) setupHandlers() {
	if c.msgHandlers == nil {
		c.msgHandlers = MessageHandlers{}
		c.AddHandler("PING", &pingHandler{})
		c.AddHandler("PRIVMSG", &defaultCTCPHandler{})
		c.AddHandler("001", &authHandler{})
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

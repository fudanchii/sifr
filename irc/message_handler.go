package irc

import (
	"strings"
)

type MessageHandler func(*Message)

type MessageHandlers map[string][]MessageHandler

// CTCP message quoted with \001
func ctcpQuote(cmd string) string {
	quoted := "\001" + cmd + "\001"
	return quoted
}

// Remove CTCP \001 quote,
// return the original string if it has no quote.
func ctcpDequote(cmd string) string {
	if cmd[0] != '\001' || cmd[len(cmd)-1] != '\001' {
		return cmd
	}
	return cmd[1 : len(cmd)-1]
}

// Add default handler for incoming messages.
func (c *Client) setupMsgHandlers() {
	c.msgHandlers = make(MessageHandlers)
	c.AddHandler("PING", func(msg *Message) {
		c.Pong(msg.Body)
	})
	c.AddHandler("PRIVMSG", c.privMsgDefaultHandler)
}

// Stolen from https://github.com/thoj/go-ircevent/blob/master/irc_callback.go
// Assign MessageHandler to Client's msgHandlers list.
// MessageHandler then will executed upon incoming Message,
// like a transparent filter chain.
func (c *Client) AddHandler(cmd string, fn MessageHandler) {
	cmd = strings.ToUpper(cmd)
	if _, ok := c.msgHandlers[cmd]; ok {
		c.msgHandlers[cmd] = append(c.msgHandlers[cmd], fn)
	} else {
		c.msgHandlers[cmd] = make([]MessageHandler, 1)
		c.msgHandlers[cmd][0] = fn
	}
}

// Default handler for PRIVMSG messages.
func (c *Client) privMsgDefaultHandler(msg *Message) {
	if msg.IsCTCP() && c.User.IsMsgForMe(msg) {
		c.handleCTCP(msg)
		return
	}
	//---- handle the usual PRIVMSG here
}

// Handle CTCP messages
// XXX: Custom CTCP handler?
func (c *Client) handleCTCP(msg *Message) {
	cmd := ctcpDequote(msg.Body)
	sender := msg.FromNick()
	switch {
	case cmd == "VERSION":
		c.ResponseCTCP(sender, "VERSION Sifr:0.0.0")
	case cmd == "SOURCE":
		c.ResponseCTCP(sender, "SOURCE https://github.com/fudanchii/sifr")
	case strings.HasPrefix(cmd, "PING "):
		c.ResponseCTCP(sender, cmd)
	}
}

// Check if this Message is CTCP message.
func (m *Message) IsCTCP() bool {
	return m.Body[0] == '\001' && m.Body[len(m.Body)-1] == '\001'
}

// Return the nick whose this Message came from.
func (m *Message) FromNick() string {
	offset := strings.Index(m.From, "!")
	if offset == -1 {
		offset = len(m.From)
	}
	return m.From[:offset]
}

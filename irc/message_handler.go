package irc

import (
	"strings"
)

type Message struct {
	// Ident whose this message coming from.
	From string

	// Ident, nick, or channel where this message sent to.
	To string

	// Purpose of the message, eg. NOTICE, or PRIVMSG, etc.
	Action string

	// Action params, separated by space.
	Params string

	// Message's body (trail).
	Body string
}

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
// TODO: Handle CTCP PING
func (c *Client) handleCTCP(msg *Message) {
	cmd := ctcpDequote(msg.Body)
	switch cmd {
	case "VERSION":
		c.responseCTCP(msg.From, "VERSION Sifr:0.0.0")
	case "SOURCE":
		c.responseCTCP(msg.From, "SOURCE https://github.com/fudanchii/sifr")
	}
}

// Check if this Message is CTCP message.
func (m *Message) IsCTCP() bool {
	return m.Body[0] == '\001' && m.Body[len(m.Body)-1] == '\001'
}

// Return the nick whose this Message came from.
func (m *Message) FromNick() string {
	return m.From[:strings.Index(m.From, "!")]
}

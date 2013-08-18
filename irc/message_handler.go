package irc

import "strings"

type Message struct {
	from   string
	to     string
	action string
	body   string
}

type MessageHandlers map[string][]func(*Message)

func ctcpQuote(cmd string) string {
	quoted := "\001" + cmd + "\001"
	return quoted
}

func ctcpDequote(cmd string) string {
	if cmd[0] != '\001' || cmd[len(cmd)-1] != '\001' {
		return cmd
	}
	return cmd[1 : len(cmd)-1]
}

func (c *Client) setupMsgHandlers() {
	c.msgHandlers = make(MessageHandlers)

	c.AddHandler("PING", func(msg *Message) {
		c.Pong(msg.body)
	})

	c.AddHandler("PRIVMSG", c.privMsgDefaultHandler)
}

func (c *Client) AddHandler(cmd string, fn func(*Message)) {
	cmd = strings.ToUpper(cmd)
	if _, ok := c.msgHandlers[cmd]; ok {
		c.msgHandlers[cmd] = append(c.msgHandlers[cmd], fn)
	} else {
		c.msgHandlers[cmd] = make([]func(*Message), 1)
		c.msgHandlers[cmd][0] = fn
	}
}

func (c *Client) privMsgDefaultHandler(msg *Message) {
	if msg.isCTCP() && c.user.isMsgForMe(msg) {
		c.handleCTCP(msg)
		return
	}
	//---- handle the usual PRIVMSG here
}

func (c *Client) handleCTCP(msg *Message) {
	cmd := ctcpDequote(msg.body)
	switch cmd {
	case "VERSION":
		c.responseCTCP(msg.from, "VERSION Sifr:0.0.0")
	case "SOURCE":
		c.responseCTCP(msg.from, "SOURCE https://github.com/fudanchii/sifr")
	}
}

func (m *Message) isCTCP() bool {
	return m.body[0] == '\001' && m.body[len(m.body)-1] == '\001'
}

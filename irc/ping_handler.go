package irc

import msg "github.com/fudanchii/sifr/irc/message"

type pingHandler struct{}

func (p *pingHandler) Handle(c *Client, m *msg.Message) {
	c.Pong(m.Body)
}

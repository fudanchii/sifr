package irc

import "github.com/fudanchii/sifr/irc/message"

type authHandler struct{}

func (a *authHandler) Handle(c *Client, m *message.Message) {
	if len(c.User.password) == 0 || c.Authenticated {
		return
	}
	c.PrivMsg("nickserv", "identify "+c.User.password)
	c.Authenticated = true
}

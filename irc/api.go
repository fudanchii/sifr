package irc

import (
	msg "github.com/fudanchii/sifr/irc/message"
)

type MessageHandler interface {
	Handle(c *Client, m *msg.Message)
}

type MessageHandlers map[string][]MessageHandler

package skill

import (
    "github.com/fudanchii/sifr/irc"
)

func ActivateFor(c *Client) {
    c.AddHandler("PRIVMSG", func (msg *irc.Message) {
        flipCoin(c, msg)
    })
}

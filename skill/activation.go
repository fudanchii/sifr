package skill

import (
    "github.com/fudanchii/sifr/irc"
)

func ActivateFor(c *irc.Client) {
    c.AddHandler("PRIVMSG", func (msg *irc.Message) {
        flipCoin(c, msg)
    })
}

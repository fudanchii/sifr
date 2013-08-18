package skill

import (
    "github.com/fudanchii/sifr/irc"
    "strings"
)

func joinChannel( c *irc.Client, msg *irc.Message) {
    if nocmd(msg.Body, ".j", true) {
        return
    }
    c.Join(strings.TrimSpace(msg.Body[3:]), "")
}


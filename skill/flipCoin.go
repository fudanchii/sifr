package skill

import (
    "github.com/fudanchii/sifr/irc"
    "math/rand"
    "strings"
    "time"
)

func parse(txt string) []string {
    result := strings.Split(txt, ",")
    if len(result) < 2 {
        result = strings.Split(txt, " ")
    }
    return result
}

func process(args []string) string {
    rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
    return args[rnd.Intn(len(args))]
}

func flip( c *Client, msg *irc.Message) {
    if nocmd(msg.body) {
        return
    }
    args := parse(msg.body)
    if context := msg.to[0]; context == '#' {
        c.PrivMsg(msg.to, msg.from + ": " + process(args))
    } else {
        c.PrivMsg(msg.from, process(args))
    }
}


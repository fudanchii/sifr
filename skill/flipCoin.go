package skill

import (
    "github.com/fudanchii/sifr/irc"
    "math/rand"
    "strings"
    "time"
)

func nocmd(txt string) bool {
    return txt[:2] != ".c"
}

func parse(txt string) []string {
    result := strings.Split(txt[3:], ",")
    if len(result) < 2 {
        result = strings.Split(txt, " ")
    }
    return result
}

func process(args []string) string {
    rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
    return args[rnd.Intn(len(args))]
}

func flipCoin( c *irc.Client, msg *irc.Message) {
    if nocmd(msg.Body) {
        return
    }
    args := parse(msg.Body)
    if context := msg.To[0]; context == '#' {
        c.PrivMsg(msg.To, msg.From + ": " + process(args))
    } else {
        c.PrivMsg(msg.From, process(args))
    }
}


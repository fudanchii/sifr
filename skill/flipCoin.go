package skill

import (
	"github.com/fudanchii/sifr/irc"
	"math/rand"
	"strings"
	"time"
)

func parse(txt string) []string {
	result := strings.Split(txt[3:], ",")
	if len(result) < 2 {
		result = strings.Split(txt[3:], " ")
	}
	return result
}

func process(args []string) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return args[rnd.Intn(len(args))]
}

func flipCoin(c *irc.Client, msg *irc.Message) {
	args := parse(msg.Body)
	if context := msg.To[0]; context == '#' {
		c.PrivMsg(msg.To, msg.FromNick()+": "+process(args))
	} else {
		c.PrivMsg(msg.FromNick(), process(args))
	}
}

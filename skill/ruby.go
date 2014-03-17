package skill

import (
	"strings"
	"github.com/fudanchii/sifr/irc"
	"github.com/mitchellh/go-mruby"
)

var (
	mrb *mruby.Mrb
)

func init() {
	mrb = mruby.NewMrb()
}

func execRuby(c *irc.Client, m *irc.Message) {
	result, err := mrb.LoadString(strings.SplitN(m.Body, " ", 2)[1])
	if err != nil {
		replyWith(c, m, err.Error())
	} else {
		replyWith(c, m, "> "+result.String())
	}
}

func replyWith(c *irc.Client, msg *irc.Message, txt string) {
	if context := msg.To[0]; context == '#' {
		c.PrivMsg(msg.To, msg.FromNick()+": "+txt)
	} else {
		c.PrivMsg(msg.FromNick(), txt)
	}
}

func resetRbCtx(c *irc.Client, m *irc.Message) {
	mrb.Close()
	mrb = mruby.NewMrb()
	replyWith(c, m, "Current context has been reset")
}

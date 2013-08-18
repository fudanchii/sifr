package irc

import "strings"


func ctcpQuote(cmd string) string {
    quoted = "\01" + cmd + "\01"
    return quoted
}

func ctcpDequote(cmd string) string {
    if cmd[0] != "\01" || cmd[len(cmd) - 1] != "\01" {
        return cmd
    }
    return cmd[1:len(cmd) - 1]
}

func (c *Client) setupMsgHandlers() {
    c.msgHandlers = make(MessageHandlers)

    c.AddHandler("PING", func (msg *Message) {
      c.Pong(msg.body)
    })

    c.AddHandler("PRIVMSG", c.privMsgDefaultHandler)
}

func (c *Client) AddHandler(cmd string, fn func(*Message)) {
    cmd = strings.ToUpper(cmd)
    if _, ok := c.msgHandlers[cmd]; ok {
        c.msgHandlers[cmd] = append(c.msgHandlers[cmd], fn)
    } else {
        c.msgHandlers[cmd] = make([]func(*Message), 1)
        c.msgHandlers[cmd][0] = fn
    }
}

func (c *Client) privMsgDefaultHandler(msg *Message) {
    if msg.isCTCP && c.isMsgForMe(msg) {
        c.handleCTCP(msg)
        return
    }
    //---- handle the usual PRIVMSG here
}

func (m *Message) isCTCP() bool {
    return m.body[0] == "\01" && m.body[len(m.body)-1] == "\01"
}

func (c *Client) isMsgForMe(msg *Message) bool {
    re, _ := regexp.Compile("(^| )" + c.user.nick + "([\\W]|$)")
    return msg.to == c.user.nick || re.MatchString(msg.body)
}


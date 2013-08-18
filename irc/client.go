package irc

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type Client struct {
	conn        net.Conn
	User        User
	msgHandlers MessageHandlers
	Errorchan   chan error
}

func Connect(addr string, user User) (*Client, error) {
	client := &Client{
		User:      user,
		Errorchan: make(chan error),
	}
	client.setupMsgHandlers()
	cConn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, &Client{}
	}
	client.conn = cConn
	client.register(user)
	go client.handleInput()
	return client, nil
}

func (c *Client) Error() string {
	return "Error creating client.\n"
}

func (c *Client) Send(cmd string, a ...interface{}) {
	str := fmt.Sprintf(cmd, a...)
	c.conn.Write([]byte(str + "\r\n"))
	log.Println("out>", str)
}

func (c *Client) Join(channel, password string) {
	c.Send("JOIN %s %s", channel, password)
}

func (c *Client) Nick(nick string) {
	c.Send("NICK " + nick)
}

func (c *Client) Notice(to, msg string) {
	c.Send("NOTICE %s :%s", to, msg)
}

func (c *Client) Part(channel string) {
	c.Send("PART " + channel)
}

func (c *Client) Ping(arg string) {
	c.Send("PING :" + arg)
}

func (c *Client) Pong(arg string) {
	c.Send("PONG :" + arg)
}

func (c *Client) PrivMsg(to, msg string) {
	c.Send("PRIVMSG %s :%s", to, msg)
}

func (c *Client) register(user User) {
	c.Nick(user.Nick)
	c.Send("USER %s %d * :%s", user.Nick, user.mode, user.Realname)
}

func (c *Client) responseCTCP(to, answer string) {
	c.Notice(to, ctcpQuote(answer))
}

func (c *Client) respondTo(maskedUser, action, talkedTo, message string) {
    maskedUser = strings.TrimPrefix(maskedUser, ":")
    message = strings.TrimPrefix(message, ":")
	user := strings.SplitN(maskedUser, "!", 2)
	msg := &Message{
		From:   user[0],
		To:     talkedTo,
		Action: action,
		Body:   message,
	}
	for _, fn := range c.msgHandlers[action] {
		go fn(msg)
	}
}

func (c *Client) handleInput() {
	defer c.conn.Close()
	reader := bufio.NewReader(c.conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			c.Errorchan <- err
			break
		}
        // FIXME: This is obviously not the right way to parse messages
		packet := strings.SplitN(line[:len(line)-2], " ", 4)
        if len(packet) == 2 {
            packet = []string{packet[1], packet[0], c.User.Nick, packet[1]}
        }
		if len(packet) == 4 {
			go c.respondTo(packet[0], packet[1], packet[2], packet[3])
		}
		log.Println("in>", line)
	}

}

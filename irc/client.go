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
	messagechan chan *Message
}

func Connect(addr string, user User) (*Client, error) {
	client := &Client{
		User:      user,
		Errorchan: make(chan error),
		// Buffer 25 messages at channel, this is a lot!
		messagechan: make(chan *Message, 25),
	}
	client.setupMsgHandlers()
	cConn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, &Client{}
	}
	client.conn = cConn
	go client.handleInput()
	go client.processMessage()
	client.register(user)
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
	if len(user.password) != 0 {
		c.PrivMsg("nickserv", "identify "+user.password)
	}
}

func (c *Client) responseCTCP(to, answer string) {
	c.Notice(to, ctcpQuote(answer))
}

func (c *Client) respondTo(maskedUser, action, talkedTo, message string) {
	message = strings.TrimPrefix(message, ":")
	maskedUser = strings.TrimPrefix(maskedUser, ":")
	user := strings.SplitN(maskedUser, "!", 2)
	msg := &Message{
		From:   user[0],
		To:     talkedTo,
		Action: action,
		Body:   message,
	}
	c.messagechan <- msg
}

func (c *Client) handleInput() {
	defer c.conn.Close()
	reader := bufio.NewReader(c.conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			close(c.messagechan)
			c.Errorchan <- err
			break
		}
		log.Println("in>", line)
		// FIXME: This is obviously not the right way to parse messages
		packet := strings.SplitN(line[:len(line)-2], " ", 4)
		if len(packet) == 2 {
			packet = []string{packet[1], packet[0], c.User.Nick, packet[1]}
		}
		if len(packet) == 4 {
			c.respondTo(packet[0], packet[1], packet[2], packet[3])
		}
	}
}

func (c *Client) processMessage() {
	for {
		msg, ok := <-c.messagechan
		if !ok {
			return
		}
		for _, fn := range c.msgHandlers[msg.Action] {
			fn(msg)
		}
	}
}

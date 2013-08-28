package irc

// IRC client implementation.

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type Client struct {
	User      User
	Errorchan chan error

	// Hold the actual irc connection
	conn net.Conn

	// List of MessageHandler chain, keyed by Message's action
	msgHandlers MessageHandlers

	// Message gets transmitted through this channel
	messagechan chan *Message
	error       error
}

// Connect to irc server at `addr` as this `user`
// if success Connect returns `Client`.
func Connect(addr string, user User) (*Client, error) {
	client := &Client{
		User:        user,
		Errorchan:   make(chan error),
		messagechan: make(chan *Message, 25),
	}
	client.setupMsgHandlers()
	if cConn, err := net.Dial("tcp", addr); err != nil {
		return nil, &Client{error: err}
	}
	client.conn = cConn
	go client.handleInput()
	go client.processMessage()
	client.register(user)
	return client, nil
}

// Implement Error interface
func (c *Client) Error() string {
	return "Error creating client: " + c.error.Error() + "\n"
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

// Register User to the server, and optionally identify with nickserv
// XXX: Need to wait nickserv identify until User actually connected.
//      - At the first CTCP VERSION request?
func (c *Client) register(user User) {
	c.Nick(user.Nick)
	c.Send("USER %s %d * :%s", user.Nick, user.mode, user.Realname)
	if len(user.password) != 0 {
		c.PrivMsg("nickserv", "identify "+user.password)
	}
}

// Response CTCP message.
func (c *Client) responseCTCP(to, answer string) {
	c.Notice(to, ctcpQuote(answer))
}

// Sit still wait for input, then pass it to Client.messagechan
func (c *Client) handleInput() {
	defer c.conn.Close()
	reader := bufio.NewReader(c.conn)
	for {
		if line, err := reader.ReadString('\n'); err != nil {
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
			c.messagechan <- createMessage(packet[0], packet[1], packet[2], packet[3])
		}
	}
}

// Execute MessageHandler chain once its arrived at Client.messagechan
func (c *Client) processMessage() {
	for {
		if msg, ok := <-c.messagechan; !ok {
			return
		}
		for _, fn := range c.msgHandlers[msg.Action] {
			fn(msg)
		}
	}
}

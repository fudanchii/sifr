package irc

// IRC client implementation.

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

type Client struct {
	User       User
	Errorchan  chan error
	Registered bool

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
		Registered:  false,
		messagechan: make(chan *Message, 25),
	}
	client.setupMsgHandlers()
	cConn, err := net.Dial("tcp", addr)
	if err != nil {
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
	if c.Registered {
		return
	}

	c.Nick(user.Nick)
	c.Send("USER %s %d * :%s", user.Nick, user.mode, user.Realname)

	if len(user.password) != 0 {
		// Sleep until we sure it's connected
		time.Sleep(time.Duration(5000) * time.Millisecond)

		c.PrivMsg("nickserv", "identify "+user.password)
	}
}

// Response CTCP message.
func (c *Client) ResponseCTCP(to, answer string) {
	c.Notice(to, ctcpQuote(answer))
}

// Sit still wait for input, then pass it to Client.messagechan
func (c *Client) handleInput() {
	defer c.conn.Close()
	scanner := bufio.NewScanner(c.conn)
	for {
		if scanner.Scan() {
			msg := scanner.Text()
			log.Println("in>", msg)
			c.messagechan <- parseMessage(msg)
		} else {
			close(c.messagechan)
			c.Errorchan <- scanner.Err()
			break
		}
	}
}

// Execute MessageHandler chain once its arrived at Client.messagechan
func (c *Client) processMessage() {
	for msg := range c.messagechan {
		for _, fn := range c.msgHandlers[msg.Action] {
			fn(msg)
		}
	}
}

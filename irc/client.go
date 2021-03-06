// IRC client implementation.
package irc

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"

	msg "github.com/fudanchii/sifr/irc/message"
)

// Connect to irc server at `addr` as this `user`
// if success Connect returns `Client`.
func Connect(addr string, user User) (*Client, error) {
	cConn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, &Client{error: err}
	}
	client := newClient(&user, cConn)
	return client, nil
}

// ConnectTLS connects to irc server via TLS/SSL
func ConnectTLS(addr string, user User, cfg *tls.Config) (*Client, error) {
	cConn, err := tls.Dial("tcp", addr, cfg)
	if err != nil {
		return nil, &Client{error: err}
	}
	client := newClient(&user, cConn)
	return client, nil
}

func newClient(user *User, conn net.Conn) *Client {
	client := &Client{
		User:        *user,
		Errorchan:   make(chan error),
		messagechan: make(chan *msg.Message, 25),
	}
	client.conn = conn
	client.setupHandlers()
	go client.handleInput()
	go client.processMessage()
	client.initReg(*user)
	return client
}

// ----------------------------------------------------------

type Client struct {
	User          User
	Errorchan     chan error
	Authenticated bool
	conn          net.Conn
	msgHandlers   MessageHandlers
	messagechan   chan *msg.Message
	error         error
}

// Error implements Error interface
func (c *Client) Error() string {
	return "Error creating client: " + c.error.Error() + "\n"
}

// Send cmd to irc.
func (c *Client) Send(cmd string, a ...interface{}) {
	str := fmt.Sprintf(cmd, a...)
	c.conn.Write([]byte(str + "\r\n"))
	log.Println("out>", str)
}

// Join to irc channel
func (c *Client) Join(channel, password string) {
	c.Send("JOIN %s %s", channel, password)
}

// Nick sets client handle name
func (c *Client) Nick(nick string) {
	c.Send("NICK " + nick)
}

// Notice sends notice message to irc
func (c *Client) Notice(to, msg string) {
	c.Send("NOTICE %s :%s", to, msg)
}

// Part tells client to leave channel
func (c *Client) Part(channel string) {
	c.Send("PART " + channel)
}

// Ping sends ping message to irc
func (c *Client) Ping(arg string) {
	c.Send("PING :" + arg)
}

// Pong replies ping message from irc
func (c *Client) Pong(arg string) {
	c.Send("PONG :" + arg)
}

// PrivMsg sends privmsg to irc
func (c *Client) PrivMsg(to, msg string) {
	c.Send("PRIVMSG %s :%s", to, msg)
}

// Response CTCP message.
func (c *Client) ResponseCTCP(to, answer string) {
	c.Notice(to, msg.TagCTCP(answer))
}

// Register User to the server, and optionally identify with nickserv
func (c *Client) initReg(user User) {
	c.Nick(user.Nick)
	c.Send("USER %s %d * :%s", user.Nick, user.mode, user.Realname)
}

// Sit still wait for input, then pass it to Client.messagechan
func (c *Client) handleInput() {
	defer c.conn.Close()
	scanner := bufio.NewScanner(c.conn)
	for {
		if scanner.Scan() {
			line := scanner.Text()
			log.Println("in>", line)
			c.messagechan <- msg.Parse(line)
		} else {
			close(c.messagechan)
			c.Errorchan <- scanner.Err()
			break
		}
	}
}

// Execute MessageHandler chain once its arrived at Client.messagechan
func (c *Client) processMessage() {
	for m := range c.messagechan {
		for _, fn := range c.msgHandlers[m.Action] {
			fn.Handle(c, m)
		}
	}
}

package irc

import (
    "bufio"
    "net"
    "strings"
)

type Client struct {
    conn    net.Conn
}

type User struct {
    nick        string
    username    string
    host        string
    realname    string
    // FIXME: Unencrypted!!!!
    password    string
}

func Connect(addr string, port uint32, user User) (*Client, error) {
    cConn, err := net.Dial("tcp", addr)
    if err != nil {
        return nil, Client{}
    }
    client = &Client{
        conn: cConn
    }
    err = client.register(user)
    if err != nil {
    }
    go client.handleInput()
    return client, nil
}

func (c *Client) Send(cmd string) {
    c.conn.Write(cmd + "\r\n")
}

func (c *Client) register(user User) {
    c.Send("NICK " + user.nick)
    c.Send("USER ")
}

func (c *Client) Pong(arg string) {
    c.Send("PONG " + arg)
}

func (c *Client) Ping() {
    c.Send("PING ")
}

func (c *Client) handleInput() {
    defer c.conn.Close()
    reader := bufio.NewReader(c.conn)
    msg := ""
    for {
        msg = reader.ReadString("\r")
    }
}

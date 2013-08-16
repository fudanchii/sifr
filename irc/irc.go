package irc

import (
    "bufio"
    "net"
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
        return nil, 
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
    c.writer.WriteString(cmd + "\r\n")
    c.writer.Flush()
}

func (c *Client) handleInput() {
    reader := bufio.NewReader(c.conn)
    defer c.conn.Close()
    for {

    }
}

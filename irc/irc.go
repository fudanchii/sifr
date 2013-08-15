package irc

import (
    "bufio"
    "net"
)

type Client struct {
    conn    net.Conn
}

func Connect(addr string, port uint32, user User) (*Client, error) {
    conn, err := net.Dial(addr, port)
    client = &Client{
        srvAddr: addr
    }
    go handleInput(*client.conn)
    return client, nil
}

func (c *Client) Send(cmd string) {
    c.writer.WriteString(cmd + "\r\n")
    c.writer.Flush()
}

func handleInput(conn *net.Conn) () {

}

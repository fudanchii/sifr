package irc

type MessageHandlers map[string][]func(*Event)

type Client struct {
    conn        net.Conn
    user        User
    msgHandlers MessageHandlers
    Errorchan   chan error
}

type User struct {
    nick        string
    mode        int
    username    string
    realname    string
    // FIXME: Unencrypted!!!!
    password    string
}

type Message struct {
    from        string
    to          string
    action      string
    body        string
}


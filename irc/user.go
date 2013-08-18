package irc

import "regexp"

type User struct {
    nick        string
    mode        int
    username    string
    realname    string
    // FIXME: Unencrypted!!!!
    password    string
}

func NewUser(nick, username, realname, password string) *User {
    return &User{
        nick:       nick
        mode:       0
        username:   username
        realname:   realname
        password:   password
    }
}

func (u *User) isMsgForMe(msg *Message) bool {
    re, _ := regexp.Compile("(^| )" + u.nick + "([\\W]|$)")
    return msg.to == u.nick || re.MatchString(msg.body)
}


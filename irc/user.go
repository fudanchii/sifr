package irc

import "regexp"

type User struct {
	Nick     string
	mode     int
	Username string
	Realname string
	// FIXME: Unencrypted!!!!
	password string
}

func NewUser(nick, username, realname, password string) *User {
	return &User{
		Nick:     nick,
		mode:     0,
		Username: username,
		Realname: realname,
		password: password,
	}
}

func (u *User) isMsgForMe(msg *Message) bool {
	re, _ := regexp.Compile("(^| )" + u.Nick + "([\\W]|$)")
	return msg.To == u.Nick || re.MatchString(msg.Body)
}

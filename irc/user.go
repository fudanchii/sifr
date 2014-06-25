package irc

import (
	"regexp"

	msg "github.com/fudanchii/sifr/irc/message"
)

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

// Check if message directed to this User, or to the channel.
func (u *User) OwnThis(m *msg.Message) bool {
	re, _ := regexp.Compile("(^| )" + u.Nick + "([\\W]|$)")
	return m.To == u.Nick || re.MatchString(m.Body)
}

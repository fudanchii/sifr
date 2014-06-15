// message package implements irc message parser, etc.
package message

import (
	"strings"
)

// IsCTCP check if this Message is CTCP message.
func (m *Message) IsCTCP() bool {
	return m.Body[0] == '\001' && m.Body[len(m.Body)-1] == '\001'
}

// FromNick returns the nick whose this Message came from.
func (m *Message) FromNick() string {
	offset := strings.Index(m.From, "!")
	if offset == -1 {
		offset = len(m.From)
	}
	return m.From[:offset]
}

// -----------------------------------------------

// TagCTCP tags message `cmd` with CTCP delimiter and
// escape unsafe characters.
func TagCTCP(cmd string) string {
	quoted := "\001" + Quote(cmd) + "\001"
	return quoted
}

// UntagCTCP removes CTCP delimiter tag and unescape
// quoted characters.
func UntagCTCP(cmd string) string {
	if cmd[0] != '\001' || cmd[len(cmd)-1] != '\001' {
		return Dequote(cmd)
	}
	return Dequote(cmd[1 : len(cmd)-1])
}

// Dequote unescapes quoted characters from cmd.
func Dequote(cmd string) string {
	str := strings.Replace(lowDequote(cmd), `\a`, "\001", -1)
	return cpyExclude(str, 0134)
}

// Quote escapes unsafe characters drom cmd.
func Quote(cmd string) string {
	str := strings.Replace(ctcpQuote(cmd), "\020", "\020\020", -1)
	str = strings.Replace(str, "\r", "\020r", -1)
	str = strings.Replace(str, "\n", "\020n", -1)
	return strings.Replace(str, "\000", "\0200", -1)
}

func ctcpQuote(cmd string) string {
	str := strings.Replace(cmd, `\`, `\\`, -1)
	return strings.Replace(str, "\001", `\a`, -1)
}

func lowDequote(cmd string) string {
	str := strings.Replace(cmd, "\0200", "\000", -1)
	str = strings.Replace(str, "\020n", "\n", -1)
	str = strings.Replace(str, "\020r", "\r", -1)
	return cpyExclude(str, 020)
}

func cpyExclude(str string, chr byte) string {
	var tch = []byte{}
	for i := 0; i < len(str); i++ {
		if str[i] == chr {
			switch {
			case i+1 == len(str):
				continue
			case str[i+1] == chr:
				i++
			default:
				continue
			}
		}
		tch = append(tch, str[i])
	}
	return string(tch)
}

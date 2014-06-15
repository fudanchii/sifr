package message

import (
	"bytes"
	"strings"
)

type Message struct {
	// Ident whose this message coming from.
	From string

	// Ident, nick, or channel where this message sent to.
	To string

	// Purpose of the message, eg. NOTICE, or PRIVMSG, etc.
	Action string

	// Action params, separated by space.
	Params string

	// Message's body (trail).
	Body string
}

func Parse(message string) *Message {
	msgStruct := &Message{}
	message = parsePrefix(message, msgStruct)
	return msgStruct
}

func doSplit(message string, holder *string) string {
	splitted := strings.SplitN(message, " ", 2)
	*holder = splitted[0]
	if len(splitted) == 2 {
		return splitted[1]
	}
	return ""
}

func parsePrefix(message string, msgStruct *Message) string {
	if message[0] == ':' { // Is this really a prefix?
		message = doSplit(message[1:], &msgStruct.From)
	}
	return parseCommand(message, msgStruct)
}

func parseCommand(message string, msgStruct *Message) string {
	message = doSplit(message, &msgStruct.Action)
	return parseTrailer(message, msgStruct)
}

func parseTrailer(message string, msgStruct *Message) string {
	if len(message) == 0 {
		return ""
	} else if message[0] == ':' { // We've got trailer here!
		msgStruct.Body = message[1:]
		return ""
	}
	return parseParams(message, msgStruct)
}

func parseParams(message string, msgStruct *Message) string {
	if len(msgStruct.To) == 0 {
		message = doSplit(message, &msgStruct.To)
	} else {
		var tempStr string
		buffer := bytes.NewBufferString(msgStruct.Params)
		message = doSplit(message, &tempStr)
		if len(msgStruct.Params) > 0 {
			buffer.WriteString(" ")
		}
		buffer.WriteString(tempStr)
		msgStruct.Params = buffer.String()
	}
	return parseTrailer(message, msgStruct)
}

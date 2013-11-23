package irc

import (
	"strings"
)

func parseMessage(message string) *Message {
	msgStruct := &Message{
		From:   "",
		To:     "",
		Action: "",
		Params: "",
		Body:   "",
	}
	message = parsePrefix(message, msgStruct)
	return msgStruct
}

func doSplit(message string, holder *string) string {
	splitted := strings.SplitN(message, " ", 2)
	holder = splitted[0]
	return splitted[1]
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
	if message[0] == ':' { // We've got trailer here!
		msgStruct.Body = message[1:]
		return ""
	}
	return parseParams(message.msgStruct)
}

func parseParams(message string, msgStruct *Message) string {
	if len(msgStruct.To) == 0 {
		message = doSplit(message, &msgStruct.To)
	} else {
		var tempStr string
		buffer := bytes.NewBufferString(msgStruct.Params)
		message = doSplit(message, &tempStr)
		buffer.WriteString(" ")
		buffer.WriteString(tempStr)
		msgStruct.Params = buffer.String()
	}
	return parseTrailer(message, msgStruct)
}

package message

import "testing"

func TestParse(t *testing.T) {
	message := &Message{
		From:   "",
		Action: "PING",
		To:     "",
		Params: "",
		Body:   "irc.example.com",
	}
	parsed := Parse("PING :irc.example.com")
	if compare(message, parsed) {
		t.Fatalf("%s == %s | %s == %s", message.Action, parsed.Action, message.Body, parsed.Body)
	}

	message = &Message{
		From:   "irc.example.com",
		Action: "439",
		To:     "*",
		Params: "",
		Body:   "Please wait while we process your connection.",
	}
	parsed = Parse(":irc.example.com 439 * :Please wait while we process your connection.")
	if compare(message, parsed) {
		t.Fatalf("%s == %s | %s == %s", message.Action, parsed.Action, message.Body, parsed.Body)
	}

	message = &Message{
		From:   "irc.example.com",
		Action: "004",
		To:     "shifuru",
		Params: "irc.example.com hybrid-7.2.3+plexus-3.1.0(20130523_0-539) CDFGNRSUWXabcdfgijklnopqrsuwxyz BIMNORSabcehiklmnopqstvz Iabehkloqv",
		Body:   "",
	}
	parsed = Parse(":irc.example.com 004 shifuru irc.example.com hybrid-7.2.3+plexus-3.1.0(20130523_0-539) CDFGNRSUWXabcdfgijklnopqrsuwxyz BIMNORSabcehiklmnopqstvz Iabehkloqv")
	if compare(message, parsed) {
		t.Fatalf("'%s' == '%s'", message.Params, parsed.Params)
	}

}

func compare(message *Message, parsed *Message) bool {
	return !(message.From == parsed.From && message.Action == parsed.Action && message.To == parsed.To &&
		message.Params == parsed.Params && message.Body == parsed.Body)
}

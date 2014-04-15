package irc_test

import (
	. "github.com/fudanchii/sifr/irc/message"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Message", func() {
	It("Shoud parse ping message", func() {
		message := &Message{
			From:   "",
			Action: "PING",
			To:     "",
			Params: "",
			Body:   "irc.example.com",
		}
		parsed := Parse("PING :irc.example.com")
		Expect(parsed).To(Equal(message))
	})
	It("Should parse any message", func() {
		message := &Message{
			From:   "irc.example.com",
			Action: "439",
			To:     "*",
			Params: "",
			Body:   "Please wait while we process your connection.",
		}
		parsed := Parse(":irc.example.com 439 * :Please wait while we process your connection.")
		Expect(parsed).To(Equal(message))
	})
	It("Should parse message with long params", func() {
		message := &Message{
			From:   "irc.example.com",
			Action: "004",
			To:     "shifuru",
			Params: "irc.example.com hybrid-7.2.3+plexus-3.1.0(20130523_0-539) CDFGNRSUWXabcdfgijklnopqrsuwxyz BIMNORSabcehiklmnopqstvz Iabehkloqv",
			Body:   "",
		}
		parsed := Parse(":irc.example.com 004 shifuru irc.example.com hybrid-7.2.3+plexus-3.1.0(20130523_0-539) CDFGNRSUWXabcdfgijklnopqrsuwxyz BIMNORSabcehiklmnopqstvz Iabehkloqv")
		Expect(parsed).To(Equal(message))
	})
	It("Should parse CTCP message", func() {
		msg := "SED \n\t\big\020\001\000\\:"
		parsed := Parse(":actor PRIVMSG victim :\001SED \020n\t\big\020\020\\a\0200\\\\:\001")
		Expect(UntagCTCP(parsed.Body)).To(Equal(msg))
		Expect(TagCTCP(msg)).To(Equal(parsed.Body))
	})
})

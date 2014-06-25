package irc_spec

import (
	. "github.com/fudanchii/sifr/irc"
	. "github.com/fudanchii/sifr/irc/message"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IRC", func() {
	var (
		user *User
	)

	BeforeEach(func() {
		user = NewUser("testnick", "testuser", "testname", "testpassword")
	})

	Describe("Current User", func() {
		It("should have nick", func() {
			Expect(user.Nick).To(Equal("testnick"))
		})
		It("should have username", func() {
			Expect(user.Username).To(Equal("testuser"))
		})
		It("should have realname", func() {
			Expect(user.Realname).To(Equal("testname"))
		})
		It("can identify if given message is received for herself", func() {
			msg := &Message{
				From:   "",
				To:     "testnick",
				Body:   "",
				Action: "",
				Params: "",
			}
			Expect(user.OwnThis(msg)).To(BeTrue())
		})
	})
})

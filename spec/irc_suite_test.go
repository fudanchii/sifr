package irc_spec

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestIRCSpec(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IRC Spec Suite")
}

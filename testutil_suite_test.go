package testutil

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTestutil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testutil Suite")
}

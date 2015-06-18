package noaareader_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestNoaareader(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Noaareader Suite")
}

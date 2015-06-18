package lab_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLab(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lab Suite")
}

package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFlashflood(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Flashflood Suite")
}

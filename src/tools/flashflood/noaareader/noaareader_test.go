package noaareader_test

import (
	"tools/flashflood/noaareader"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"tools/flashflood/messagewriter"
	"github.com/cloudfoundry/sonde-go/events"
	"time"
)

var _ = Describe("NoaaReader", func() {

	It("consumer reads messages from traffic controller with valid auth token", func(done Done) {
		defer close(done)
		messageChan := make(chan *events.LogMessage, 1)
		writer := messagewriter.NewChannelLogWriter(messageChan)
		consumer := noaareader.NewFakeConsumer(messageChan)
		reader := noaareader.NewNoaaReader("app-guid", "good token", consumer)

		writer.Send(1, time.Now())

		logMessage, ok := reader.Read()
		Expect(ok).To(BeTrue())
		Expect(logMessage.GetMessage()).To(ContainSubstring("1"))
	})

})

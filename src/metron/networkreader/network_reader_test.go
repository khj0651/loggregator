package networkreader_test

import (
	"net"

	"metron/networkreader"
	"metron/writers/mocks"

	"github.com/cloudfoundry/loggregatorlib/loggertesthelper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NetworkReader", func() {
	Context("with a reader running", func() {
		var reader *networkreader.NetworkReader
		var writer mocks.MockByteArrayWriter

		BeforeEach(func() {
			writer = mocks.MockByteArrayWriter{}
			reader = networkreader.New("127.0.0.1:3456", "networkReader", &writer, loggertesthelper.Logger())

			loggertesthelper.TestLoggerSink.Clear()
			go reader.Start()

			Eventually(loggertesthelper.TestLoggerSink.LogContents).Should(ContainSubstring("Listening on port 127.0.0.1:3456"))
		})

		AfterEach(func() {
			reader.Stop()
		})

		It("sends data recieved on UDP socket to its writer", func() {
			expectedData := "Some Data"
			otherData := "More stuff"

			connection, err := net.Dial("udp", "localhost:3456")

			_, err = connection.Write([]byte(expectedData))
			Expect(err).NotTo(HaveOccurred())

			Eventually(writer.Data).Should(HaveLen(1))
			data := string(writer.Data()[0])
			Expect(data).To(Equal(expectedData))

			_, err = connection.Write([]byte(otherData))
			Expect(err).NotTo(HaveOccurred())

			Eventually(writer.Data).Should(HaveLen(2))

			data = string(writer.Data()[1])
			Expect(data).To(Equal(otherData))
		})
	})
})

package new_integration_test

import (
	"net"
	"time"

	"github.com/cloudfoundry/storeadapter"
	"github.com/gogo/protobuf/proto"


	"github.com/cloudfoundry/dropsonde/dropsonde_unmarshaller"
	"github.com/cloudfoundry/loggregatorlib/agentlistener"
	"github.com/cloudfoundry/loggregatorlib/loggertesthelper"
	"github.com/cloudfoundry/sonde-go/events"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "integration_tests/metron/matchers"
)

var _ = FDescribe("Self Instrumentation", func() {
	var (
		fakeDoppler agentlistener.AgentListener
		envelopes   chan *events.Envelope
	)

	BeforeEach(func() {
		envelopes = make(chan *events.Envelope, 1024)

		var signedMessages <-chan []byte
		unsignedMessages := make(chan []byte)
		fakeDoppler, signedMessages = agentlistener.NewAgentListener("localhost:3457", loggertesthelper.Logger(), "fakeDoppler")
		unmarshaller := dropsonde_unmarshaller.NewDropsondeUnmarshaller(loggertesthelper.Logger())

		go unmarshaller.Run(unsignedMessages, envelopes)
		go func() {
			for signedMessage := range signedMessages {
				unsignedMessages <- signedMessage[32:]
			}
		}()

		go fakeDoppler.Start()

		announceToEtcd()
	})

	AfterEach(func() {
		fakeDoppler.Stop()
	})

	It("sends metrics about the Dropsonde network reader", func() {
		metronInput, _ := net.Dial("udp", "localhost:51161")

		metronInput.Write(basicValueMessage())
		expected := events.Envelope{
			Origin:    proto.String("MetronAgent"),
			EventType: events.Envelope_CounterEvent.Enum(),
			CounterEvent: &events.CounterEvent{
				Name:  proto.String("receivedMessageCount"),
				Delta: proto.Uint64(1),
				Total: proto.Uint64(1),
			},
		}

		Eventually(envelopes).Should(Receive(MatchSpecifiedContents(&expected)))
	})
})

func announceToEtcd() {
	node := storeadapter.StoreNode{
		Key:   "/healthstatus/doppler/z1/0",
		Value: []byte("localhost"),
	}

	adapter := etcdRunner.Adapter()
	adapter.Create(node)
	adapter.Disconnect()
	time.Sleep(50 * time.Millisecond)
}

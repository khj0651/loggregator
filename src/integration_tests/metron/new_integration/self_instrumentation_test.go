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

	. "integration_tests/metron/matchers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Self Instrumentation", func() {
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
				Name:  proto.String("dropsondeAgentListener.receivedMessageCount"),
				Delta: proto.Uint64(1),
				Total: proto.Uint64(1),
			},
		}

		Eventually(envelopes).Should(Receive(MatchSpecifiedContents(&expected)))
	})

	It("counts legacy unmarshal errors", func() {
		metronInput, _ := net.Dial("udp", "localhost:51160")
		metronInput.Write([]byte{1, 2, 3})

		expected := events.Envelope{
			Origin:    proto.String("MetronAgent"),
			EventType: events.Envelope_CounterEvent.Enum(),
			CounterEvent: &events.CounterEvent{
				Name:  proto.String("legacyUnmarshaller.unmarshalErrors"),
				Delta: proto.Uint64(1),
				Total: proto.Uint64(1),
			},
		}

		Eventually(envelopes).Should(Receive(MatchSpecifiedContents(&expected)))
	})

	Describe("for Dropsonde unmarshaller", func() {
		It("counts errors", func() {
			metronInput, _ := net.Dial("udp", "localhost:51161")
			metronInput.Write([]byte{1, 2, 3})

			expected := events.Envelope{
				Origin:    proto.String("MetronAgent"),
				EventType: events.Envelope_CounterEvent.Enum(),
				CounterEvent: &events.CounterEvent{
					Name:  proto.String("dropsondeUnmarshaller.unmarshalErrors"),
					Delta: proto.Uint64(1),
					Total: proto.Uint64(1),
				},
			}

			Eventually(envelopes).Should(Receive(MatchSpecifiedContents(&expected)))
		})

		It("counts unmarshalled Dropsonde messages by type", func() {
			metronInput, _ := net.Dial("udp", "localhost:51161")
			metronInput.Write(basicValueMessage())

			expected := events.Envelope{
				Origin:    proto.String("MetronAgent"),
				EventType: events.Envelope_CounterEvent.Enum(),
				CounterEvent: &events.CounterEvent{
					Name:  proto.String("dropsondeUnmarshaller.valueMetricReceived"),
					Delta: proto.Uint64(1),
					Total: proto.Uint64(1),
				},
			}

			Eventually(envelopes).Should(Receive(MatchSpecifiedContents(&expected)))
		})

		It("counts log messages specially", func() {
			metronInput, _ := net.Dial("udp", "localhost:51161")

			logEnvelope := &events.Envelope{
				Origin:    proto.String("fake-origin-2"),
				EventType: events.Envelope_LogMessage.Enum(),
				LogMessage: &events.LogMessage{
					Message:     []byte("hello"),
					MessageType: events.LogMessage_OUT.Enum(),
					Timestamp:   proto.Int64(1234),
				},
			}
			logBytes, _ := proto.Marshal(logEnvelope)

			metronInput.Write(logBytes)

			expected := events.Envelope{
				Origin:    proto.String("MetronAgent"),
				EventType: events.Envelope_CounterEvent.Enum(),
				CounterEvent: &events.CounterEvent{
					Name:  proto.String("dropsondeUnmarshaller.logMessageTotal"),
					Delta: proto.Uint64(1),
					Total: proto.Uint64(1),
				},
			}

			Eventually(envelopes).Should(Receive(MatchSpecifiedContents(&expected)))
		})

		It("counts unknown event types", func() {
			metronInput, _ := net.Dial("udp", "localhost:51161")
			message := basicValueMessageEnvelope()
			message.EventType = events.Envelope_EventType(2000).Enum()
			bytes, err := proto.Marshal(message)
			Expect(err).ToNot(HaveOccurred())

			metronInput.Write(bytes)

			message = basicValueMessageEnvelope()
			badEventType := events.Envelope_EventType(1000)
			message.EventType = &badEventType
			bytes, err = proto.Marshal(message)
			Expect(err).ToNot(HaveOccurred())

			metronInput.Write(bytes)

			expected := events.Envelope{
				Origin:    proto.String("MetronAgent"),
				EventType: events.Envelope_CounterEvent.Enum(),
				CounterEvent: &events.CounterEvent{
					Name:  proto.String("dropsondeUnmarshaller.unknownEventTypeReceived"),
					Delta: proto.Uint64(2),
					Total: proto.Uint64(2),
				},
			}

			Eventually(envelopes).Should(Receive(MatchSpecifiedContents(&expected)))
		})
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

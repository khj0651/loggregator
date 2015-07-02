package metrics_test

import (
	"net"
	"time"

	"github.com/cloudfoundry/storeadapter"
	"github.com/gogo/protobuf/proto"

	"github.com/cloudfoundry/sonde-go/events"

	. "integration_tests/metron/matchers"

	"integration_tests/metron/metrics"

	"metron/writers/messageaggregator"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Self Instrumentation", func() {
	var (
		testDoppler *metrics.TestDoppler
	)

	BeforeEach(func() {
		testDoppler = metrics.NewTestDoppler()
		go testDoppler.Start()

		announceToEtcd()
	})

	AfterEach(func() {
		testDoppler.Stop()
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

		Eventually(testDoppler.MessageChan, 2).Should(Receive(MatchSpecifiedContents(&expected)))
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

		Eventually(testDoppler.MessageChan).Should(Receive(MatchSpecifiedContents(&expected)))
	})

	Describe("for Message Aggregator", func() {
		var (
			metronInput net.Conn
			originalTTL time.Duration
		)

		BeforeEach(func() {
			metronInput, _ = net.Dial("udp", "localhost:51161")
			originalTTL = messageaggregator.MaxTTL
		})

		AfterEach(func() {
			metronInput.Close()
			messageaggregator.MaxTTL = originalTTL
		})

		It("emits metrics for counter events", func() {
			metronInput.Write(basicCounterEvent())
			metronInput.Write(basicCounterEvent())

			expected := events.Envelope{
				Origin:    proto.String("MetronAgent"),
				EventType: events.Envelope_CounterEvent.Enum(),
				CounterEvent: &events.CounterEvent{
					Name:  proto.String("MessageAggregator.counterEventReceived"),
					Delta: proto.Uint64(1),
					Total: proto.Uint64(2),
				},
			}

			Eventually(testDoppler.MessageChan).Should(Receive(MatchSpecifiedContents(&expected)))

			expected = events.Envelope{
				Origin:    proto.String("MetronAgent"),
				EventType: events.Envelope_CounterEvent.Enum(),
				CounterEvent: &events.CounterEvent{
					Name:  proto.String("MessageAggregator.counterEventReceived"),
					Delta: proto.Uint64(1),
					Total: proto.Uint64(3),
				},
			}
			Consistently(testDoppler.MessageChan).ShouldNot(Receive(MatchSpecifiedContents(&expected)))
		})

		It("emits metrics for http start", func() {
			metronInput.Write(basicHTTPStartEvent())

			expected := events.Envelope{
				Origin:    proto.String("MetronAgent"),
				EventType: events.Envelope_CounterEvent.Enum(),
				CounterEvent: &events.CounterEvent{
					Name:  proto.String("MessageAggregator.httpStartReceived"),
					Delta: proto.Uint64(1),
					Total: proto.Uint64(1),
				},
			}

			Eventually(testDoppler.MessageChan).Should(Receive(MatchSpecifiedContents(&expected)))
		})

		It("emits metrics for http stop", func() {
			metronInput.Write(basicHTTPStopEvent())

			expected := events.Envelope{
				Origin:    proto.String("MetronAgent"),
				EventType: events.Envelope_CounterEvent.Enum(),
				CounterEvent: &events.CounterEvent{
					Name:  proto.String("MessageAggregator.httpStopReceived"),
					Delta: proto.Uint64(1),
					Total: proto.Uint64(1),
				},
			}

			Eventually(testDoppler.MessageChan).Should(Receive(MatchSpecifiedContents(&expected)))
		})

		It("emits metrics for unmatched http stop", func() {
			metronInput.Write(basicHTTPStopEvent())

			expected := events.Envelope{
				Origin:    proto.String("MetronAgent"),
				EventType: events.Envelope_CounterEvent.Enum(),
				CounterEvent: &events.CounterEvent{
					Name:  proto.String("MessageAggregator.httpUnmatchedStopReceived"),
					Delta: proto.Uint64(1),
					Total: proto.Uint64(1),
				},
			}

			Eventually(testDoppler.MessageChan).Should(Receive(MatchSpecifiedContents(&expected)))
		})

		It("emits metrics for http start stop", func() {
			metronInput.Write(basicHTTPStartEvent())
			metronInput.Write(basicHTTPStopEvent())

			expected := events.Envelope{
				Origin:    proto.String("MetronAgent"),
				EventType: events.Envelope_CounterEvent.Enum(),
				CounterEvent: &events.CounterEvent{
					Name:  proto.String("MessageAggregator.httpStartStopEmitted"),
					Delta: proto.Uint64(1),
					Total: proto.Uint64(1),
				},
			}

			Eventually(testDoppler.MessageChan, 2).Should(Receive(MatchSpecifiedContents(&expected)))
		})

		It("emits metrics for uncategorized events", func() {
			message := basicValueMessageEnvelope()
			message.EventType = events.Envelope_LogMessage.Enum()
			bytes, err := proto.Marshal(message)
			Expect(err).ToNot(HaveOccurred())

			metronInput.Write(bytes)

			expected := events.Envelope{
				Origin:    proto.String("MetronAgent"),
				EventType: events.Envelope_CounterEvent.Enum(),
				CounterEvent: &events.CounterEvent{
					Name:  proto.String("MessageAggregator.uncategorizedEvents"),
					Delta: proto.Uint64(1),
					Total: proto.Uint64(1),
				},
			}

			Eventually(testDoppler.MessageChan, 2).Should(Receive(MatchSpecifiedContents(&expected)))
		})
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

			Eventually(testDoppler.MessageChan).Should(Receive(MatchSpecifiedContents(&expected)))
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

			Eventually(testDoppler.MessageChan).Should(Receive(MatchSpecifiedContents(&expected)))
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

			Eventually(testDoppler.MessageChan).Should(Receive(MatchSpecifiedContents(&expected)))
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

			Eventually(testDoppler.MessageChan).Should(Receive(MatchSpecifiedContents(&expected)))
		})

		It("does not forward unknown events", func() {
			metronInput, _ := net.Dial("udp", "localhost:51161")
			message := basicValueMessageEnvelope()
			message.EventType = events.Envelope_EventType(2000).Enum()
			bytes, err := proto.Marshal(message)
			Expect(err).ToNot(HaveOccurred())

			metronInput.Write(bytes)

			Eventually(testDoppler.MessageChan).ShouldNot(Receive(MatchSpecifiedContents(&message)))

			message = basicValueMessageEnvelope()
			badEventType := events.Envelope_EventType(1000)
			message.EventType = &badEventType
			bytes, err = proto.Marshal(message)
			Expect(err).ToNot(HaveOccurred())

			metronInput.Write(bytes)

			Consistently(testDoppler.MessageChan).ShouldNot(Receive(MatchSpecifiedContents(&message)))
		})
	})

	Describe("for Dropsonde marshaller", func() {
		It("counts marshalled Dropsonde messages by type", func() {
			metronInput, _ := net.Dial("udp", "localhost:51161")
			metronInput.Write(basicValueMessage())

			expected := events.Envelope{
				Origin:    proto.String("MetronAgent"),
				EventType: events.Envelope_CounterEvent.Enum(),
				CounterEvent: &events.CounterEvent{
					Name:  proto.String("dropsondeMarshaller.valueMetricMarshalled"),
					Delta: proto.Uint64(1),
					Total: proto.Uint64(1),
				},
			}

			Eventually(testDoppler.MessageChan, 2).Should(Receive(MatchSpecifiedContents(&expected)))
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

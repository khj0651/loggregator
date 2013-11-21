package sinkserver

import (
	"github.com/cloudfoundry/gosteno"
	"github.com/cloudfoundry/loggregatorlib/cfcomponent/instrumentation"
	"github.com/cloudfoundry/loggregatorlib/loggertesthelper"
	"github.com/cloudfoundry/loggregatorlib/logmessage"
	messagetesthelpers "github.com/cloudfoundry/loggregatorlib/logmessage/testhelpers"
	"github.com/stretchr/testify/assert"
	"loggregator/groupedsinks"
	"loggregator/sinks"
	"testing"
	"time"
)

type testSink struct {
	channel chan *logmessage.Message
}

func (ts testSink) Emit() instrumentation.Context {
	return instrumentation.Context{}
}

func (ts testSink) AppId() string {
	return "appId"
}

func (ts testSink) Run() {

}

func (ts testSink) Channel() chan *logmessage.Message {

	return ts.channel
}

func (ts testSink) Identifier() string {
	return "testSink"
}

func (ts testSink) Logger() *gosteno.Logger {
	return loggertesthelper.Logger()
}

func TestDumpToSinkWithLessThan20Messages(t *testing.T) {
	testMessageRouter := NewMessageRouter(1024, loggertesthelper.Logger())

	activeSinks := groupedsinks.NewGroupedSinks()
	dumpSink := sinks.NewDumpSink("appId", 100, loggertesthelper.Logger())
	activeSinks.Register(dumpSink)

	message := messagetesthelpers.NewMessage(t, "message 1", "appId")
	for i := 0; i < 19; i++ {
		dumpSink.Channel() <- message
	}
	close(dumpSink.Channel())
	<-time.After(10 * time.Millisecond)

	sink := testSink{make(chan *logmessage.Message, 100)}
	testMessageRouter.dumpToSink(sink, activeSinks)

	assert.Equal(t, 19, len(sink.Channel()))
}

func TestDumpToSinkLimitsMessagesTo20(t *testing.T) {
	testMessageRouter := NewMessageRouter(1024, loggertesthelper.Logger())
	sink := testSink{make(chan *logmessage.Message, 100)}
	activeSinks := groupedsinks.NewGroupedSinks()
	dumpSink := sinks.NewDumpSink("appId", 100, loggertesthelper.Logger())

	message := messagetesthelpers.NewMessage(t, "message 1", "appId")
	for i := 0; i < 100; i++ {
		dumpSink.Channel() <- message
	}

	activeSinks.Register(dumpSink)
	testMessageRouter.dumpToSink(sink, activeSinks)

	assert.Equal(t, 20, len(sink.Channel()))

}
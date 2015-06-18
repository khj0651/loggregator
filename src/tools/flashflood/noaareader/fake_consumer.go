package noaareader
import (
	"github.com/cloudfoundry/sonde-go/events"
)

type fakeConsumer struct {
	messageChan chan *events.LogMessage
}

func NewFakeConsumer(messageChan chan *events.LogMessage) *fakeConsumer{
	return &fakeConsumer {
		messageChan : messageChan,
	}
}

func (f *fakeConsumer) TailingLogsWithoutReconnect(appGuid string, authToken string, outputChan chan<- *events.LogMessage) error {
	
	for message := range f.messageChan {
		outputChan <- message
	}

	return nil
}
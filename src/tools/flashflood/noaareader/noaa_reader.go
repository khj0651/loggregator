package noaareader
import (
	"github.com/cloudfoundry/sonde-go/events"
)

type NoaaConsumer interface {
	TailingLogsWithoutReconnect(appGuid string, authToken string, outputChan chan<- *events.LogMessage) error
}

type noaaReader struct {
	consumer NoaaConsumer
	messageChan chan *events.LogMessage
}

func NewNoaaReader(appGuid string, authToken string, consumer NoaaConsumer) *noaaReader {
	messageChan := make(chan *events.LogMessage)
	go consumer.TailingLogsWithoutReconnect(appGuid, authToken, messageChan)
	return &noaaReader{
		consumer: consumer,
		messageChan: messageChan,
	}
}

func (nr *noaaReader) Read() (*events.LogMessage, bool) {
	logMessage, ok := <- nr.messageChan
	return logMessage, ok
}



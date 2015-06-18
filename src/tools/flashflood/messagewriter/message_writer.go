package messagewriter

import (
	"io"
	"time"
)

type messageWriter struct {
	w io.Writer
	*roundSequencer
}

func NewMessageWriter(w io.Writer) *messageWriter {
	return &messageWriter{
		w:              w,
		roundSequencer: newRoundSequener(),
	}
}

func (mw *messageWriter) Send(roundId uint, timestamp time.Time) {
	mw.w.Write(formatMsg(roundId, mw.nextSequence(roundId), timestamp))
}

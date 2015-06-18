package server

import (
	"fmt"
	"github.com/cloudfoundry/sonde-go/events"
	"net/http"
	"os"
	"tools/flashflood/experiment"
	"tools/flashflood/lab"
	"tools/flashflood/messagereader"
	"tools/flashflood/messagewriter"
)

type Server struct {
	l *lab.Lab
}

func New(lr messagereader.LoggregatorReader) *Server {
	return &Server{
		l: lab.New(func() lab.Experiment {
			writer := messagewriter.NewMessageWriter(os.Stdout)
			reader := messagereader.NewMessageReader(lr)
			return experiment.New(writer, reader, experiment.NewRound)
		}),
	}
}

func (s Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("/start", func(w http.ResponseWriter, req *http.Request) {
		s.l.Start()
		fmt.Fprintf(w, `Starting the flood. Visit <a href="/results">results</a> to see what's happening.`)
	})

	mux.HandleFunc("/stop", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Stopping the flood … soonish.")

		s.l.Stop()
	})

	mux.HandleFunc("/results", func(w http.ResponseWriter, req *http.Request) {
		s.l.Results().Render(w)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// The "/" pattern matches everything, so we need to check
		// that we're at the root here.
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
	})

	mux.ServeHTTP(w, req)
}

func tailLogs() <-chan *events.LogMessage {
	//	noaa.NewConsumer(tcURL, )
	return make(chan *events.LogMessage)
}

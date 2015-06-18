package server_test

import (
	"tools/flashflood/server"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"github.com/cloudfoundry/gunk/urljoiner"
	"net/http"
	"strings"
	"regexp"
	"tools/flashflood/messagereader"
	"github.com/cloudfoundry/sonde-go/events"
)

var _ = Describe("Server", func() {
	var testServer *httptest.Server

	BeforeEach(func() {
		c:= make(chan *events.LogMessage)
		testServer = httptest.NewServer(server.New(messagereader.NewChannelReader(c)))
	})

	AfterEach(func() {
		testServer.Close()
	})

	It("responds to the proper endpoints", func() {
		resp, err := http.Get(urljoiner.Join(testServer.URL, "/start"))
		Expect(err).NotTo(HaveOccurred())

		body, err := ioutil.ReadAll(resp.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(body).To(ContainSubstring(`Starting the flood. Visit <a href="/results">results</a> to see what's happening.`))

		resp, err = http.Get(urljoiner.Join(testServer.URL, "/stop"))
		Expect(err).NotTo(HaveOccurred())

		body, err = ioutil.ReadAll(resp.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(body).To(ContainSubstring("Stopping the flood … soonish."))
	})

	It("displays no results when '/start' has not been hit", func() {
		resp, err := http.Get(urljoiner.Join(testServer.URL, "/results"))
		Expect(err).NotTo(HaveOccurred())

		body, err := ioutil.ReadAll(resp.Body)
		Expect(err).NotTo(HaveOccurred())

		Expect(body).To(ContainSubstring("There are no results. Please issue a request to the /start endpoint"))
		Expect(body).NotTo(ContainSubstring("Rate"))
		Expect(body).NotTo(ContainSubstring("Messages Sent"))
		Expect(body).NotTo(ContainSubstring("Messages Received"))
		Expect(body).NotTo(ContainSubstring("Message Loss Percentage"))
	})

	It("puts results when '/start' has been hit", func() {
		_, err := http.Get(urljoiner.Join(testServer.URL, "/start"))
		Expect(err).NotTo(HaveOccurred())


		getResults := func() string {
			resp, err := http.Get(urljoiner.Join(testServer.URL, "/results"))
			Expect(err).NotTo(HaveOccurred())

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			return compactHTML(string(body))
		}


		htmlResult := compactHTML(
		`<tr>
			<th>Message Rate</th>
			<th>Messages Sent</th>
			<th>Messages Received</th>
			<th>Message Loss Percentage</th>
		</tr>`)
		Eventually(getResults).Should(ContainSubstring(htmlResult))
	})
})

func compactHTML(formattedHTML string) string {
	result := strings.Replace(formattedHTML, "\n", "", -1)
	re := regexp.MustCompile(">\\s+<")
	result = re.ReplaceAllString(result, "><")

	return result
}

